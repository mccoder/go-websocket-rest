go get -u github.com/mailru/easyjson/...

DROP2_BASEDIR=$(dirname "$0")
(
    cd $DROP2_BASEDIR/model
    easyjson -all WsRequestMessage.go
    easyjson -all WsResponceMessage.go
    easyjson -all WsSubscribeAction.go
)