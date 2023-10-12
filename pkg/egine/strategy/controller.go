package strategy

type Controller struct {
	strategy  Strategy
	dataframe *model.Dataframe
	broker    service.Broker
	started   bool
}
