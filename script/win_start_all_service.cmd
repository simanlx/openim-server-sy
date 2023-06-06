SET ROOT=%cd%
cd %ROOT%\..\bin\
start cmd /C .\crazy_server_api.exe -port 10002
start cmd /C .\crazy_server_cms_api.exe -port 10006
start cmd /C .\crazy_server_user.exe -port 10110
start cmd /C .\crazy_server_friend.exe -port 10120
start cmd /C .\crazy_server_group.exe -port 10150
start cmd /C .\crazy_server_auth.exe -port 10160
start cmd /C .\crazy_server_admin_cms.exe -port 10200
start cmd /C .\crazy_server_message_cms.exe -port 10190
start cmd /C .\crazy_server_statistics.exe -port 10180
start cmd /C .\crazy_server_msg.exe -port 10130
start cmd /C .\crazy_server_office.exe -port 10210
start cmd /C .\crazy_server_organization.exe -port 10220
start cmd /C .\crazy_server_conversation.exe -port 10230
start cmd /C .\crazy_server_cache.exe -port 10240
start cmd /C .\crazy_server_push.exe -port 10170
start cmd /C .\crazy_server_msg_transfer.exe
start cmd /C .\crazy_server_sdk_server.exe -openIM_api_port 10002 -openIM_ws_port 10001 -sdk_ws_port 10003 -openIM_log_level 6
start cmd /C .\crazy_server_msg_gateway.exe -rpc_port 10140 -ws_port 10001
start cmd /C .\crazy_server_demo.exe -port 10004
cd %ROOT%