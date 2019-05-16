package kafka

import (
	"bufio"
)

const (
	groupCoordinator       = 0
	transactionCoordinator = 1
)

// FindCoordinatorRequestV0 requests the coordinator for the specified group or transaction
//
// See http://kafka.apache.org/protocol.html#The_Messages_FindCoordinator
type findCoordinatorRequestV0 struct {
	// CoordinatorKey holds id to use for finding the coordinator (for groups, this is
	// the groupId, for transactional producers, this is the transactional id)
	CoordinatorKey string
}

type findCoordinatorRequestV1 struct {
	// CoordinatorKey holds id to use for finding the coordinator (for groups, this is
	// the groupId, for transactional producers, this is the transactional id)
	CoordinatorKey string
	// The type of coordinator to find (0 = group, 1 = transaction)
	CoordinatorType int8
}

func (t findCoordinatorRequestV0) size() int32 {
	return sizeofString(t.CoordinatorKey)
}

func (t findCoordinatorRequestV1) size() int32 {
	return sizeofString(t.CoordinatorKey) + sizeof(t.CoordinatorType)
}

func (t findCoordinatorRequestV0) writeTo(w *bufio.Writer) {
	writeString(w, t.CoordinatorKey)
}

func (t findCoordinatorRequestV1) writeTo(w *bufio.Writer) {
	writeString(w, t.CoordinatorKey)
	writeInt8(w, t.CoordinatorType)
}

type findCoordinatorResponseCoordinatorV0 struct {
	// NodeID holds the broker id.
	NodeID int32

	// Host of the broker
	Host string

	// Port on which broker accepts requests
	Port int32
}

func (t findCoordinatorResponseCoordinatorV0) size() int32 {
	return sizeofInt32(t.NodeID) +
		sizeofString(t.Host) +
		sizeofInt32(t.Port)
}

func (t findCoordinatorResponseCoordinatorV0) writeTo(w *bufio.Writer) {
	writeInt32(w, t.NodeID)
	writeString(w, t.Host)
	writeInt32(w, t.Port)
}

func (t *findCoordinatorResponseCoordinatorV0) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt32(r, size, &t.NodeID); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.Host); err != nil {
		return
	}
	if remain, err = readInt32(r, remain, &t.Port); err != nil {
		return
	}
	return
}

type findCoordinatorResponseV0 struct {
	// ErrorCode holds response error code
	ErrorCode int16

	// Coordinator holds host and port information for the coordinator
	Coordinator findCoordinatorResponseCoordinatorV0
}

type findCoordinatorResponseV1 struct {
	// Duration in milliseconds for which the request was throttled due to
	// quota violation (Zero if the request did not violate any quota)
	ThrottleTimeMs int32
	// Response error code
	ErrorCode int16
	// Response error message
	ErrorMessage string
	// Host and port information for the coordinator
	Coordinator findCoordinatorResponseCoordinatorV0
}

func (t findCoordinatorResponseV0) size() int32 {
	return sizeofInt16(t.ErrorCode) +
		t.Coordinator.size()
}

func (t findCoordinatorResponseV1) size() int32 {
	return sizeof(t.ThrottleTimeMs) + sizeof(t.ErrorCode) +
		sizeof(t.ErrorMessage) + t.Coordinator.size()
}

func (t findCoordinatorResponseV0) writeTo(w *bufio.Writer) {
	writeInt16(w, t.ErrorCode)
	t.Coordinator.writeTo(w)
}

func (t findCoordinatorResponseV1) writeTo(w *bufio.Writer) {
	writeInt32(w, t.ThrottleTimeMs)
	writeInt16(w, t.ErrorCode)
	writeString(w, t.ErrorMessage)
	t.Coordinator.writeTo(w)
}

func (t *findCoordinatorResponseV0) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt16(r, size, &t.ErrorCode); err != nil {
		return
	}
	if remain, err = (&t.Coordinator).readFrom(r, remain); err != nil {
		return
	}
	return
}

func (t *findCoordinatorResponseV1) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt32(r, size, &t.ThrottleTimeMs); err != nil {
		return
	}
	if remain, err = readInt16(r, remain, &t.ErrorCode); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.ErrorMessage); err != nil {
		return
	}
	if remain, err = (&t.Coordinator).readFrom(r, remain); err != nil {
		return
	}
	return
}
