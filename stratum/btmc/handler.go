package btmc

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	"github.com/bytom/btmcpool/common/logger"
	ss "github.com/bytom/btmcpool/stratum"
	"github.com/bytom/btmcpool/stratum/btmc/util"
	"github.com/bytom/consensus/difficulty"
	"github.com/bytom/protocol/bc"
)

type encodeMethod uint

const (
	spHandleLogin   encodeMethod = 0
	spHandleGetWork encodeMethod = 1
	spHandleSubmit  encodeMethod = 2
)

func (request *LoginReq) handleLogin(session *ss.TcpSession) error {
	// 1. check session state
	if session.GetState() != ss.SStateConnected {
		session.Error(request.Id, ss.ErrorMultipleAuth)
		return fmt.Errorf("wrong session state %v", session.GetState())
	}

	// 2. create and register worker
	worker, err := ss.NewWorker(request.Params.Login, "")
	if err != nil {
		session.Error(request.Id, ss.ErrorFormatAuthorize)
		return err
	}

	miner := worker.GetFullName()
	if err := session.GetServerState().GetConnCtrl().JudgeMiner(miner, session.SessionCtl); err != nil {
		return err
	}

	session.GetSessionData().SetWorker(worker)

	// 3. session is now authorized
	session.SetState(ss.SStateAuthorized)

	// 4. send new job in session
	if err := spJobReply(spHandleLogin, request, session); err != nil {
		return err
	}
	logger.Info("handle login",
		"session_id", session.GetId(),
		"session_ip", session.GetIp(),
		"miner", worker.GetFullName(),
		"method", request.Method)
	return nil
}

func (request *SubmitReq) handleSubmit(session *ss.TcpSession) error {
	// ensure tcp session state is authorized
	if session.GetState() != ss.SStateAuthorized {
		session.Error(request.Id, ss.ErrorUnauthorized)
		return fmt.Errorf("wrong session state %v", session.GetState())
	}

	// check if the worker is registered
	worker := session.GetSessionData().GetWorker()
	if worker == nil || worker.GetFullName() != request.Params.Id {
		session.Error(request.Id, ss.ErrorUnauthorized)
		return ss.ErrNotRegistered
	}
	miner := worker.GetFullName()

	if err := session.GetServerState().GetConnCtrl().JudgeMiner(miner, session.SessionCtl); err != nil {
		return err
	}

	// check if the job exists
	jobId, err := ss.StringToJobId(request.Params.JobId)
	if err != nil {
		session.SessionCtl.MinerErrCnt++
		session.Error(request.Id, ss.ErrorFormatSubmit)
		logger.Warn("invalid job id",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId)
		return nil
	}
	job := session.FindJob(jobId)
	if job == nil {
		session.SessionCtl.MinerErrCnt++
		session.Error(request.Id, ss.ErrorJobNotFound)
		logger.Warn("invalid job",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId)
		return nil
	}

	btmcJob := job.(*btmcJob)

	// check nonce form
	noncePattern, _ := regexp.Compile("^[0-9a-f]{1,16}$")
	if !noncePattern.MatchString(request.Params.Nonce) {
		session.SessionCtl.MinerErrCnt++
		session.Error(request.Id, ss.ErrorFormatSubmit)
		logger.Warn("invalid nonce format",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId,
			"nonce", request.Params.Nonce,
			"height", btmcJob.height)
		return nil
	}

	nonce, err := strconv.ParseUint(request.Params.Nonce, 16, 64)
	if err != nil {
		session.SessionCtl.MinerErrCnt++
		session.Error(request.Id, ss.ErrorFormatSubmit)
		logger.Warn("invalid nonce value",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId,
			"nonce", request.Params.Nonce,
			"height", btmcJob.height)
		return nil
	}

	netBits := difficulty.CompactToBig(btmcJob.bits)
	netDiff := big.NewInt(0).Div(util.GetDividend(), netBits)
	shareDiff := btmcJob.diff.Uint64()

	// check stale share
	btmcBlockTemplate := session.GetServerState().GetBlockTemplate().(*btmcBlockTemplate)
	if btmcJob.height != btmcBlockTemplate.height {
		session.SessionCtl.MinerErrCnt++
		session.Error(request.Id, ss.ErrorFormatSubmit)
		// session.SendJob()
		logger.Warn("stale share",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId,
			"nonce", request.Params.Nonce,
			"share_diff", shareDiff,
			"net_diff", netDiff,
			"height", btmcJob.height)
		return nil
	}

	// create share and verify it
	share := &btmcShare{
		nonce:     nonce,
		result:    request.Params.Result,
		job:       btmcJob,
		worker:    worker,
		netDiff:   netDiff,
		blockHash: &bc.Hash{},
	}
	if err := session.GetVerifier().Verify(share); err != nil {
		session.Error(request.Id, ss.ErrorFormatSubmit)
		logger.Warn("failed verification",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId,
			"nonce", request.Params.Nonce,
			"share_diff", shareDiff,
			"net_diff", netDiff,
			"height", btmcJob.height,
			"error", err)
		return nil
	}

	// send reply according to verification result
	switch share.GetState() {
	case ss.ShareStateAccepted:
		session.SessionCtl.MinerAcCnt++
		if err := session.Reply(request.Id, &statusReply{Status: "OK"}); err != nil {
			logger.Error("failed to send reply",
				"method", request.Method,
				"session_id", session.GetId(),
				"session_ip", session.GetIp(),
				"miner", miner,
				"job_id", request.Params.JobId,
				"height", btmcJob.height,
				"error", err)
		}
	case ss.ShareStateBlock:
		// 8.1 submit block to node
		session.SessionCtl.MinerAcCnt++
		logger.Info("found block",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId,
			"nonce", nonce,
			"share_diff", shareDiff,
			"net_diff", share.netDiff.Uint64(),
			"height", share.header.Height,
			"hash", share.blockHash.String())
		go func() {
			// TODO: need Immediately refresh current BT after submit to node
			if err := session.GetNodeSyncer().Submit(share); err != nil {
				logger.Error("failed to submit block",
					"method", request.Method,
					"session_id", session.GetId(),
					"session_ip", session.GetIp(),
					"miner", miner,
					"job_id", request.Params.JobId,
					"height", btmcJob.height,
					"error", err)
			}
		}()

		if err := session.Reply(request.Id, &statusReply{Status: "OK"}); err != nil {
			logger.Error("failed to send reply",
				"method", request.Method,
				"session_id", session.GetId(),
				"session_ip", session.GetIp(),
				"miner", miner,
				"job_id", request.Params.JobId,
				"height", btmcJob.height,
				"error", err)
		}
	case ss.ShareStateRejected:
		session.SessionCtl.MinerErrCnt++
		logger.Info("failed share",
			"session_id", session.GetId(),
			"session_ip", session.GetIp(),
			"miner", miner,
			"job_id", request.Params.JobId,
			"nonce", nonce,
			"share_diff", shareDiff,
			"net_diff", share.netDiff.Uint64(),
			"hash", share.blockHash.String(),
			"height", btmcJob.height,
			"reason", share.reason.String())
		if err := session.Error(request.Id, share.GetReason().Error()); err != nil {
			logger.Error("failed to send reply",
				"method", request.Method,
				"session_id", session.GetId(),
				"session_ip", session.GetIp(),
				"miner", miner,
				"job_id", request.Params.JobId,
				"height", btmcJob.height,
				"error", err)
		}
	}

	return nil
}

func spJobReply(method encodeMethod, request interface{}, session *ss.TcpSession) error {
	job, err := session.GetJob()
	if err != nil {
		return err
	}

	var replyErr interface{}
	switch method {
	case spHandleLogin:
		nRequest := request.(*LoginReq)
		msg := job.(*btmcJob).encodeLogin(nRequest.Params.Login)
		replyErr = session.Reply(nRequest.Id, msg)
	case spHandleGetWork:
		nRequest := request.(*GetWorkReq)
		msg := job.(*btmcJob).genReplyData()
		replyErr = session.Reply(nRequest.Id, msg)
	}
	if replyErr != nil {
		return err
	}
	if err := session.AddJob(job); err != nil {
		return err
	}
	message, err := job.Encode()
	if err != nil {
		return err
	}
	return session.Notify(message)
}
