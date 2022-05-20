package hazelcast

//
//import (
//	"github.com/hazelcast/hazelcast-go-client/serialization"
//)
//
//const (
//	dummyClassID              = 1
//	dummyFactoryID            = 1
//	s1               = "state"
//	s2            = "typeName"
//
//)
//
//type dummyPortable struct {
//	Dummy
//	logs []*pb.WorkflowLogRecord
//}
//
//// ClassID returns workflow class id for Hazelcast portable.
//func (*dummyPortable) ClassID() int32 {
//	return dummyClassID
//}
//
//// FactoryID returns workflow factory id for Hazelcast portable.
//func (*dummyPortable) FactoryID() int32 {
//	return dummyFactoryID
//}
//
//type workflowFactory struct {
//}
//
//func (*workflowFactory) Create(classID int32) serialization.Portable {
//	if classID == workflowClassID {
//		return &workflowPortable{}
//	}
//	return nil
//}
//
//func (*workflowFactory) FactoryID() int32 {
//	return workflowFactoryID
//}
//
//// NewPortableFactory returns Portable factory for Workflow metadata struct.
//func NewPortableFactory() serialization.PortableFactory {
//	return &workflowFactory{}
//}
//
//// WritePortable serializes workflow struct as Hazelcast portable.
//func (w *workflowPortable) WritePortable(writer serialization.PortableWriter) {
//	binStartOp, err := proto.MarshalOptions{Deterministic: true}.Marshal(w.StartOperation)
//	if err != nil {
//		// we can panic here, this will result in a hzerrors.ErrHazelcastSerialization
//		panic(err)
//	}
//
//	writer.WriteByte(stateFieldName, w.State[0])
//	writer.WriteString(typeNameFieldName, w.TypeName)
//	writer.WriteInt32(iterationIDFieldName, w.IterationId)
//	writer.WriteInt32(ttlFieldName, w.Ttl)
//	writer.WriteInt64(finalTTLFieldName, w.FinalTtl)
//	writer.WriteInt32(maxRetriesFieldName, w.MaxRetries)
//	writer.WriteInt32(nextActivitiesFieldName, w.NextActivity)
//	writer.WriteInt32(maxActivitiesFieldName, w.MaxActivities)
//	writer.WriteInt64(expiresAtFieldName, w.ExpiresAt)
//	writer.WriteByteArray(startOperationFieldName, binStartOp)
//	writer.WriteBool(retryableFieldName, w.Retryable)
//	writer.WriteBool(rollbackableFieldName, w.Rollbackable)
//	writer.WriteBool(replayingResultsFieldName, w.ReplayingResults)
//	writer.WriteBool(operationInProgressFieldName, w.OperationInProgress)
//	writer.WriteByteArray(traceIDFieldName, w.TraceId)
//	writer.WriteByteArray(rootSpanIDFieldName, w.RootSpanId)
//
//	out := writer.GetRawDataOutput()
//	out.WriteInt32(int32(len(w.logs)))
//	for _, l := range w.logs {
//		binLogs, err := proto.MarshalOptions{Deterministic: true}.Marshal(l)
//		if err != nil {
//			// we can panic here, this will result in a hzerrors.ErrHazelcastSerialization
//			panic(err)
//		}
//		out.WriteByteArray(binLogs)
//	}
//}
//
//// ReadPortable deserializes workflow struct as Hazelcast portable.
//func (w *workflowPortable) ReadPortable(reader serialization.PortableReader) {
//	w.Workflow = &pb.Workflow{}
//	w.State = []byte{reader.ReadByte(stateFieldName)}
//	w.TypeName = reader.ReadString(typeNameFieldName)
//	w.IterationId = reader.ReadInt32(iterationIDFieldName)
//	w.Ttl = reader.ReadInt32(ttlFieldName)
//	w.FinalTtl = reader.ReadInt64(finalTTLFieldName)
//	w.MaxRetries = reader.ReadInt32(maxRetriesFieldName)
//	w.NextActivity = reader.ReadInt32(nextActivitiesFieldName)
//	w.MaxActivities = reader.ReadInt32(maxActivitiesFieldName)
//	w.ExpiresAt = reader.ReadInt64(expiresAtFieldName)
//	binStartOp := reader.ReadByteArray(startOperationFieldName)
//	if binStartOp != nil {
//		op := &pb.HTTPRequest{}
//		if err := proto.Unmarshal(binStartOp, op); err != nil {
//			// we can panic here, this will result in a hzerrors.ErrHazelcastSerialization
//			panic(err)
//		}
//		w.StartOperation = op
//	}
//	w.Retryable = reader.ReadBool(retryableFieldName)
//	w.Rollbackable = reader.ReadBool(rollbackableFieldName)
//	w.ReplayingResults = reader.ReadBool(replayingResultsFieldName)
//	w.OperationInProgress = reader.ReadBool(operationInProgressFieldName)
//	w.TraceId = reader.ReadByteArray(traceIDFieldName)
//	w.RootSpanId = reader.ReadByteArray(rootSpanIDFieldName)
//
//	in := reader.GetRawDataInput()
//	logsCount := in.ReadInt32()
//	w.logs = make([]*pb.WorkflowLogRecord, 0, int(logsCount))
//	for i := 0; i < int(logsCount); i++ {
//		l := &pb.WorkflowLogRecord{}
//		if err := proto.Unmarshal(in.ReadByteArray(), l); err != nil {
//			// we can panic here, this will result in a hzerrors.ErrHazelcastSerialization
//			panic(err)
//		}
//		w.logs = append(w.logs, l)
//	}
//}
