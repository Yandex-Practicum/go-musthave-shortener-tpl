package logger

import "go.uber.org/zap"

var Log zap.SugaredLogger

func Initialize() error {

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}


	Log = *logger.Sugar()
	return nil
}

