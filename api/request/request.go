package request

type GetDataByUUIDRequest struct {
	UUID string `uri:"uuid" binding:"required,validUUID"`
}
