package business

type PointTypeInfo CarbonTypeInfo

func (info PointTypeInfo) CarbonTypeInfo() CarbonTypeInfo {
	return CarbonTypeInfo(info)
}
