// ==========================================================================
// 物联网快速开发自动生成dao操作代码，无需手动修改，重新生成不会自动覆盖.
// 生成日期：2022-06-27 23:24:19
// 生成路径: github.com/morgeq/iotfast/internal/app/mqtt/dao/mqtt_topic_record.go
// 生成人：dwx
// ==========================================================================

package dao

import (
	"github.com/morgeq/iotfast/internal/app/mqtt/dao/internal"
)

// internalmqttTopicRecordDao is internal type for wrapping internal DAO implements.
type internalMqttTopicRecordDao = *internal.MqttTopicRecordDao

// mqttTopicRecordDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type mqttTopicRecordDao struct {
	internalMqttTopicRecordDao
}

var (
	// MqttTopicRecord is globally public accessible object for table tools_gen_table operations.
	MqttTopicRecord = mqttTopicRecordDao{
		internal.NewMqttTopicRecordDao(),
	}
)
