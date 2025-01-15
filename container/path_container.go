package container

type BestPath struct {
	SourceIP string
	TargetIP string
}

func (bp *BestPath) Init(sourceIP string, targetIP string) {
	bp.SourceIP = sourceIP
	bp.TargetIP = targetIP
}

func (bp *BestPath) UpdateBestPath(bestPath BestPath) {
	bp.SourceIP = bestPath.SourceIP
	bp.TargetIP = bestPath.TargetIP
}
