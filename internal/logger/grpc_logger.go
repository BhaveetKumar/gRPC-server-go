package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	logIDKey     = "x-log-id"
	sessionIDKey = "x-session-id"
)

func UnaryServerInterceptor(base *Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		md, _ := metadata.FromIncomingContext(ctx)
		logID := firstOrDefault(md[logIDKey], uuid.NewString())
		sessionID := firstOrDefault(md[sessionIDKey], uuid.NewString())

		log := base.WithContext(logID, sessionID)

		if base.enableRequestIDs {
			log.Info(fmt.Sprintf("incoming request %s | input: %+v", info.FullMethod, req))
		} else {
			log.std.Println(fmt.Sprintf("INFO: incoming request %s | input: %+v", info.FullMethod, req))
		}

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			if base.enableRequestIDs {
				log.Error(fmt.Sprintf("request failed in %s | error: %v", duration.String(), err))
			} else {
				log.std.Println(fmt.Sprintf("ERROR: request failed in %s | error: %v", duration.String(), err))
			}
			return resp, err
		}

		if base.enableRequestIDs {
			log.Info(fmt.Sprintf("request succeeded in %s | output: %+v", duration.String(), resp))
		} else {
			log.std.Println(fmt.Sprintf("INFO: request succeeded in %s | output: %+v", duration.String(), resp))
		}
		return resp, nil
	}
}

func firstOrDefault(values []string, def string) string {
	if len(values) == 0 {
		return def
	}
	if values[0] == "" {
		return def
	}
	return values[0]
}
