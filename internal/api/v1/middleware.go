package minerva_api_v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	schema "github.com/Minerva-System/minerva-go/internal/schema"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

func (self *Server) TenantCheckMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		companyId := ctx.Param("company")

		log.Debug("Retrieving a tenant service worker...")
		conn, err := self.Collection.TenantSvc.Get(ctx)
		if err != nil {
			log.Error("Failed to retrieve a tenant service worker: %v", err)
			ctx.AbortWithStatusJSON(500, schema.ErrorMessage{
				Status:  500,
				Message: "Could not connect to tenant service",
			})
			return
		}

		client := rpc.NewTenantClient(conn)
		log.Debug("Checking tenant existence for %s", companyId)
		response, err := client.Exists(ctx, &rpc.EntityIndex{Index: companyId})
		if err != nil {
			log.Error("Failed to retrieve whether a tenant exists: %v", err)
			m := schema.ErrorMessage{}.FromGrpcError(err)
			m.Details = fmt.Sprintf("Failed to retrieve whether a tenant exists: %s", m.Details)
			ctx.AbortWithStatusJSON(m.Status, m)
			return
		}

		if !response.Value {
			log.Error("Failed to find company \"%s\"", companyId)
			ctx.AbortWithStatusJSON(404, schema.ErrorMessage{
				Status:  404,
				Message: "This tenant does not exist",
			})
			return
		}

		log.Debug("Tenant %s exists, continuing request", companyId)
		conn.Close()
		ctx.Next()
	}
}
