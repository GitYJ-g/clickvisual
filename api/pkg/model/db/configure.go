package db

import (
	"fmt"

	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/clickvisual/clickvisual/api/internal/invoker"
)

type Configuration struct {
	BaseModel

	K8SCmId     int    `gorm:"column:k8s_cm_id;type:int(11);" json:"k8sConfigmapId"` // config map id
	Name        string `gorm:"column:name;type:varchar(64);" json:"name"`
	Content     string `gorm:"column:content;type:longtext" json:"content"`
	Format      string `gorm:"column:format;type:varchar(32)" json:"format"`
	Version     string `gorm:"column:version;type:varchar(64)" json:"version"`
	Uid         int    `gorm:"column:uid;type:int(11) unsigned" json:"uid"`
	PublishTime int64  `gorm:"column:publish_time;type:int(11)" json:"publishTime"`
	LockUid     int    `gorm:"column:lock_uid;type:int(11) unsigned" json:"lockUid"`
	LockAt      int64  `gorm:"column:lock_at;type:bigint(11) unsigned" json:"lockAt"`

	K8SConfigMap K8SConfigMap `gorm:"foreignKey:K8SCmId;references:ID" json:"-"`
}

func (c *Configuration) TableName() string {
	return TableNameConfiguration
}

// FileName ..
func (c Configuration) FileName() string {
	return fmt.Sprintf("%s.%s", c.Name, c.Format)
}

// ConfigurationCreate CRUD
func ConfigurationCreate(db *gorm.DB, data *Configuration) (err error) {
	if err = db.Create(data).Error; err != nil {
		invoker.Logger.Error("create cluster error", zap.Error(err))
		return
	}
	return
}

// ConfigurationUpdate ...
func ConfigurationUpdate(db *gorm.DB, paramId int, ups map[string]interface{}) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}
	if err = db.Table(TableNameConfiguration).Where(sql, binds...).Updates(ups).Error; err != nil {
		invoker.Logger.Error("update cluster error", zap.Error(err))
		return
	}
	return
}

// ConfigurationInfoX Info extension method to query a single record according to Cond
func ConfigurationInfoX(conds map[string]interface{}) (resp Configuration, err error) {
	sql, binds := egorm.BuildQuery(conds)
	invoker.Logger.Debug("ConfigurationInfoX", elog.Any("conds", sql))
	if err = invoker.Db.Table(TableNameConfiguration).Unscoped().Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("K8SConfigMapInfoX infoX error", zap.Error(err))
		return
	}
	return resp, nil
}

func ConfigurationInfo(paramId int) (resp Configuration, err error) {
	var sql = "`id`= ?"
	var binds = []interface{}{paramId}
	if err = invoker.Db.Table(TableNameConfiguration).Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("cluster info error", zap.Error(err))
		return
	}
	return
}

// ConfigurationDelete 硬删除
func ConfigurationDelete(db *gorm.DB, id int) (err error) {
	if err = db.Model(Configuration{}).Delete(&Configuration{}, id).Error; err != nil {
		invoker.Logger.Error("cluster delete error", zap.Error(err))
		return
	}
	return
}

// ConfigurationList return item list by condition
func ConfigurationList(conds egorm.Conds) (resp []*Configuration, err error) {
	sql, binds := egorm.BuildQuery(conds)
	// Fetch record with Rancher Info....
	if err = invoker.Db.Table(TableNameConfiguration).Where(sql, binds...).Find(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("list clusters error", elog.String("err", err.Error()))
		return
	}
	return
}

// ConfigurationListPage return item list by pagination
func ConfigurationListPage(conds egorm.Conds, reqList *ReqPage) (total int64, respList []*Configuration) {
	respList = make([]*Configuration, 0)
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := egorm.BuildQuery(conds)
	db := invoker.Db.Table(TableNameConfiguration).Where(sql, binds...)
	db.Count(&total)
	db.Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Find(&respList)
	return
}

type ConfigurationHistory struct {
	BaseModel
	Uid             int    `gorm:"column:uid;type:int(11) unsigned" json:"uid"`
	ConfigurationId int    `gorm:"column:configuration_id;type:int(11) unsigned" json:"configurationId"`
	ChangeLog       string `gorm:"column:change_log;type:longtext" json:"changeLog"`
	Content         string `gorm:"column:content;type:longtext" json:"content"`
	Version         string `gorm:"column:version;type:varchar(64)" json:"version"`

	Configuration Configuration `json:"configuration,omitempty" gorm:"foreignKey:ConfigurationId;references:ID"`
}

func (m *ConfigurationHistory) TableName() string {
	return TableNameConfigurationHistory
}

// ConfigurationHistoryCreate CRUD
func ConfigurationHistoryCreate(db *gorm.DB, data *ConfigurationHistory) (err error) {
	if err = db.Create(data).Error; err != nil {
		invoker.Logger.Error("ConfigurationHistoryCreate cluster error", zap.Error(err))
		return
	}
	return
}

// ConfigurationHistoryUpdate ...
func ConfigurationHistoryUpdate(db *gorm.DB, paramId int, ups map[string]interface{}) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}
	if err = db.Table(TableNameConfigurationHistory).Where(sql, binds...).Updates(ups).Error; err != nil {
		invoker.Logger.Error("ConfigurationHistoryUpdate cluster error", zap.Error(err))
		return
	}
	return
}

func ConfigurationHistoryInfo(paramId int) (resp ConfigurationHistory, err error) {
	var sql = "`id`= ?"
	var binds = []interface{}{paramId}
	if err = invoker.Db.Table(TableNameConfigurationHistory).Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("ConfigurationHistoryInfo info error", zap.Error(err))
		return
	}
	return
}

// ConfigurationHistoryInfoX get single item by condition
func ConfigurationHistoryInfoX(conds map[string]interface{}) (resp ConfigurationHistory, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = invoker.Db.Table(TableNameConfigurationHistory).Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("ConfigurationHistoryInfoX infoX error", zap.Error(err))
		return
	}
	return
}

// ConfigurationHistoryDelete soft delete item by id
func ConfigurationHistoryDelete(db *gorm.DB, id int) (err error) {
	if err = db.Model(ConfigurationHistory{}).Delete(&ConfigurationHistory{}, id).Error; err != nil {
		invoker.Logger.Error("ConfigurationHistoryDelete delete error", zap.Error(err))
		return
	}
	return
}

// ConfigurationHistoryList return item list by condition
func ConfigurationHistoryList(conds egorm.Conds) (resp []*ConfigurationHistory, err error) {
	sql, binds := egorm.BuildQuery(conds)
	// Fetch record with Rancher Info....
	if err = invoker.Db.Table(TableNameConfigurationHistory).Where(sql, binds...).Find(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("ConfigurationHistoryList clusters error", elog.String("err", err.Error()))
		return
	}
	return
}

// ConfigurationHistoryListPage return item list by pagination
func ConfigurationHistoryListPage(conds egorm.Conds, reqList *ReqPage) (total int64, respList []*ConfigurationHistory) {
	respList = make([]*ConfigurationHistory, 0)
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := egorm.BuildQuery(conds)
	db := invoker.Db.Table(TableNameConfigurationHistory).Where(sql, binds...)
	db.Count(&total)
	db.Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Order("id DESC").Find(&respList)
	return
}

type ConfigurationPublish struct {
	BaseModel

	Uid                    uint `gorm:"column:uid;type:int(11) unsigned" json:"uid"`
	ConfigurationId        uint `gorm:"column:configuration_id;type:int(11) unsigned" json:"configurationId"`
	ConfigurationHistoryId uint `gorm:"column:configuration_history_id;type:int(11) unsigned" json:"configurationHistoryId"`
}

func (m *ConfigurationPublish) TableName() string {
	return TableNameConfigurationPublish
}

// ConfigurationPublishCreate CRUD
func ConfigurationPublishCreate(db *gorm.DB, data *ConfigurationPublish) (err error) {
	if err = db.Create(data).Error; err != nil {
		invoker.Logger.Error("ConfigurationPublishCreate cluster error", zap.Error(err))
		return
	}
	return
}

// ConfigurationPublishUpdate ...
func ConfigurationPublishUpdate(db *gorm.DB, paramId int, ups map[string]interface{}) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}
	if err = db.Table(TableNameConfigurationPublish).Where(sql, binds...).Updates(ups).Error; err != nil {
		invoker.Logger.Error("ConfigurationPublishUpdate cluster error", zap.Error(err))
		return
	}
	return
}

func ConfigurationPublishInfo(paramId int) (resp ConfigurationPublish, err error) {
	var sql = "`id`= ?"
	var binds = []interface{}{paramId}
	if err = invoker.Db.Table(TableNameConfigurationPublish).Where(sql, binds...).First(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("ConfigurationPublishInfo info error", zap.Error(err))
		return
	}
	return
}

// ConfigurationPublishDelete soft delete item by id
func ConfigurationPublishDelete(db *gorm.DB, id int) (err error) {
	if err = db.Model(ConfigurationPublish{}).Delete(&ConfigurationPublish{}, id).Error; err != nil {
		invoker.Logger.Error("ConfigurationPublishDelete delete error", zap.Error(err))
		return
	}
	return
}

// ConfigurationPublishList return item list by condition
func ConfigurationPublishList(conds egorm.Conds) (resp []*ConfigurationPublish, err error) {
	sql, binds := egorm.BuildQuery(conds)
	// Fetch record with Rancher Info....
	if err = invoker.Db.Table(TableNameConfigurationPublish).Where(sql, binds...).Find(&resp).Error; err != nil && err != gorm.ErrRecordNotFound {
		invoker.Logger.Error("ConfigurationPublishList clusters error", elog.String("err", err.Error()))
		return
	}
	return
}

// ConfigurationPublishListPage return item list by pagination
func ConfigurationPublishListPage(conds egorm.Conds, reqList *ReqPage) (total int64, respList []*ConfigurationPublish) {
	respList = make([]*ConfigurationPublish, 0)
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := egorm.BuildQuery(conds)
	db := invoker.Db.Table(TableNameConfigurationPublish).Where(sql, binds...)
	db.Count(&total)
	db.Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Find(&respList)
	return
}
