# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [atomix/multiraft/v1/config.proto](#atomix_multiraft_v1_config-proto)
    - [ClusterConfig](#atomix-multiraft-v1-ClusterConfig)
    - [MemberConfig](#atomix-multiraft-v1-MemberConfig)
    - [MultiRaftConfig](#atomix-multiraft-v1-MultiRaftConfig)
    - [NodeConfig](#atomix-multiraft-v1-NodeConfig)
    - [PartitionConfig](#atomix-multiraft-v1-PartitionConfig)
    - [ReplicaConfig](#atomix-multiraft-v1-ReplicaConfig)
    - [ServiceConfig](#atomix-multiraft-v1-ServiceConfig)
  
    - [MemberConfig.Role](#atomix-multiraft-v1-MemberConfig-Role)
  
- [Scalar Value Types](#scalar-value-types)



<a name="atomix_multiraft_v1_config-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## atomix/multiraft/v1/config.proto



<a name="atomix-multiraft-v1-ClusterConfig"></a>

### ClusterConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| replicas | [ReplicaConfig](#atomix-multiraft-v1-ReplicaConfig) | repeated |  |
| partitions | [PartitionConfig](#atomix-multiraft-v1-PartitionConfig) | repeated |  |






<a name="atomix-multiraft-v1-MemberConfig"></a>

### MemberConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_id | [uint64](#uint64) |  |  |
| role | [MemberConfig.Role](#atomix-multiraft-v1-MemberConfig-Role) |  |  |






<a name="atomix-multiraft-v1-MultiRaftConfig"></a>

### MultiRaftConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| heartbeat_period | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| election_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| snapshot_entry_threshold | [uint64](#uint64) |  |  |
| compaction_retain_entries | [uint64](#uint64) |  |  |
| data_dir | [string](#string) |  |  |






<a name="atomix-multiraft-v1-NodeConfig"></a>

### NodeConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_id | [uint64](#uint64) |  |  |
| host | [string](#string) |  |  |
| port | [int32](#int32) |  |  |
| config | [MultiRaftConfig](#atomix-multiraft-v1-MultiRaftConfig) |  |  |






<a name="atomix-multiraft-v1-PartitionConfig"></a>

### PartitionConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| partition_id | [uint32](#uint32) |  |  |
| host | [string](#string) |  |  |
| port | [int32](#int32) |  |  |
| members | [MemberConfig](#atomix-multiraft-v1-MemberConfig) | repeated |  |






<a name="atomix-multiraft-v1-ReplicaConfig"></a>

### ReplicaConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node_id | [uint64](#uint64) |  |  |
| host | [string](#string) |  |  |
| api_port | [int32](#int32) |  |  |
| raft_port | [int32](#int32) |  |  |






<a name="atomix-multiraft-v1-ServiceConfig"></a>

### ServiceConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| leader | [string](#string) |  |  |
| followers | [string](#string) | repeated |  |





 


<a name="atomix-multiraft-v1-MemberConfig-Role"></a>

### MemberConfig.Role


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| MEMBER | 1 |  |
| OBSERVER | 2 |  |
| WITNESS | 3 |  |


 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers ??? if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers ??? if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

