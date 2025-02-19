package permission

import (
	"strconv"

	"github.com/clickvisual/clickvisual/api/internal/service/permission"
	"github.com/clickvisual/clickvisual/api/internal/service/permission/pmsplugin"
	"github.com/clickvisual/clickvisual/api/pkg/component/core"
	"github.com/clickvisual/clickvisual/api/pkg/model/view"
)

// GetTableCurrentRoleAssignInAllDom ...
func GetTableCurrentRoleAssignInAllDom(c *core.Context) {
	var err error
	aidStr := c.Query("appId")
	// reqDom := c.Query("reqDom")
	appId, err := strconv.Atoi(aidStr)
	if err != nil {
		c.JSONE(1, "invalid appId.", err.Error())
		return
	}
	reqPerm := view.ReqPermission{
		UserId:      c.Uid(),
		ObjectType:  pmsplugin.PrefixInstance,
		ObjectIdx:   aidStr,
		SubResource: pmsplugin.Role,
		Acts:        []string{pmsplugin.ActView},
		DomainType:  pmsplugin.SystemDom,
		DomainId:    "",
	}
	if err := permission.Manager.CheckNormalPermission(reqPerm); err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}
	// currently, reqDom from fe is empty.
	// appAssignInfo := permission.Manager.GetAppRolesAssignmentInfoInDom(appId, reqDom)
	appAssignInfo := permission.Manager.GetTableRolesAssignmentInfoInAllDom(appId)
	c.JSONOK(appAssignInfo)
}

func GetTableAvailableRoles(c *core.Context) {
	aidStr := c.Query("appId")
	appId, err := strconv.Atoi(aidStr)
	if err != nil {
		c.JSONE(1, "invalid appId.", err.Error())
		return
	}
	appAvailableRoles := permission.Manager.GetTableAvailableRoles(appId)
	c.JSONOK(appAvailableRoles)
}

func ReAssignTableRolesUser(c *core.Context) {
	var err error
	reqModel := view.TableRolesAssignmentInfo{}
	err = c.Bind(&reqModel)
	if err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}
	var appId = reqModel.AppId
	reqPerm := view.ReqPermission{
		UserId:      c.Uid(),
		ObjectType:  pmsplugin.PrefixInstance,
		ObjectIdx:   strconv.Itoa(appId),
		SubResource: pmsplugin.InstanceBase,
		Acts:        []string{pmsplugin.ActEdit},
		DomainType:  pmsplugin.SystemDom,
		DomainId:    "",
	}
	if err := permission.Manager.CheckNormalPermission(reqPerm); err != nil {
		c.JSONE(1, err.Error(), nil)
		return
	}

	for i, reqRole := range reqModel.RolesInfo {
		reqModel.RolesInfo[i].BelongType = pmsplugin.PrefixInstance
		reqModel.RolesInfo[i].ReferId = appId
		if reqRole.DomainId == 0 {
			reqModel.RolesInfo[i].DomainType = ""
		}
	}
	permission.Manager.AssignTableRolesUser(reqModel.AppId, reqModel.RolesInfo)

	c.JSONOK()
}
