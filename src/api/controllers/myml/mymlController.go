package myml

import (
	mymlService "github.com/lgjafabian/MyML/src/api/services/myml"
	"github.com/emikohmann/itacademy2018-myml/src/api/util/apierrors"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func GetInformation(c *gin.Context) {

	orderID, err := strconv.ParseInt(c.Param("orderID"), 10, 64)
	if err != nil {
        c.JSON(http.StatusBadRequest, apierrors.ApiError{StatusCode: http.StatusBadRequest, Error:      err.Error(),})
		return
	}

	order, apiErr := mymlService.GetInformation(orderID)
    if apiErr != nil {
        c.JSON(apiErr.StatusCode, apiErr)
        return
	}
	
    c.JSON(http.StatusOK, order)

}