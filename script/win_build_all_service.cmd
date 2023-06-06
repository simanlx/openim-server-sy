SET ROOT=%cd%
mkdir %ROOT%\..\bin\
cd ..\cmd\crazy_server_api\ && go build -ldflags="-w -s" && move crazy_server_api.exe %ROOT%\..\bin\
cd ..\..\cmd\crazy_server_cms_api\ && go build -ldflags="-w -s" && move crazy_server_cms_api.exe %ROOT%\..\bin\
cd ..\..\cmd\crazy_server_demo\ && go build -ldflags="-w -s" && move crazy_server_demo.exe %ROOT%\..\bin\
cd ..\..\cmd\crazy_server_msg_gateway\ && go build -ldflags="-w -s" && move crazy_server_msg_gateway.exe %ROOT%\..\bin\
cd ..\..\cmd\crazy_server_msg_transfer\ && go build -ldflags="-w -s" && move crazy_server_msg_transfer.exe %ROOT%\..\bin\
cd ..\..\cmd\crazy_server_push\ && go build -ldflags="-w -s" && move crazy_server_push.exe %ROOT%\..\bin\
cd ..\..\cmd\rpc\crazy_server_admin_cms\&& go build -ldflags="-w -s" && move crazy_server_admin_cms.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_auth\&& go build -ldflags="-w -s" && move crazy_server_auth.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_cache\&& go build -ldflags="-w -s" && move crazy_server_cache.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_conversation\&& go build -ldflags="-w -s" && move crazy_server_conversation.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_friend\&& go build -ldflags="-w -s" && move crazy_server_friend.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_group\&& go build -ldflags="-w -s" && move crazy_server_group.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_message_cms\&& go build -ldflags="-w -s" && move crazy_server_message_cms.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_msg\&& go build -ldflags="-w -s" && move crazy_server_msg.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_office\&& go build -ldflags="-w -s" && move crazy_server_office.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_organization\&& go build -ldflags="-w -s" && move crazy_server_organization.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_statistics\&& go build -ldflags="-w -s" && move crazy_server_statistics.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\crazy_server_user\&& go build -ldflags="-w -s" && move crazy_server_user.exe %ROOT%\..\bin\
cd ..\..\..\cmd\Open-IM-SDK-Core\ws_wrapper\cmd\&& go build -ldflags="-w -s" crazy_server_sdk_server.go && move crazy_server_sdk_server.exe %ROOT%\..\bin\
cd %ROOT%