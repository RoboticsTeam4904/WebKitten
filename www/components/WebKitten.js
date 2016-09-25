const React = require('react')
const Immutable = require('immutable')
import {EventEmitter} from 'fbemitter'
import { Subsystem, ConnectionStatusIndicator } from './index'
import { STATUS_DISCONNECTED, STATUS_CONNECTING, STATUS_CONNECTED, STATUS_FAILED } from './ConnectionStatusIndicator'

const SUBSYSTEMS = Immutable.List.of('chassis', 'intake', 'flywheel')

const WEBSOCKET_ADDRESS = 'ws://localhost:8080/ws'
const WEBSOCKET_PROTOCOL = 'webkitten-v1'

const WEBSOCKET_RETRY_COUNT = 3

class WebKitten extends React.Component {
  constructor(props) {
    super(props)

    this.logEmitter = new EventEmitter()
    this.websocketRetriesRemaining = WEBSOCKET_RETRY_COUNT
    this.websocketHasOpened = false
  }

  componentWillMount() {
    this.openWebsocket()
  }

  componentWillDismount() {
    this.websocketConnection.close()
  }

  openWebsocket() {
    // set status to connecting
    this.setState({ connectionStatus: STATUS_CONNECTING })

    // attempt to open websocket (async)
    this.websocketConnection = new WebSocket(WEBSOCKET_ADDRESS, WEBSOCKET_PROTOCOL)

    // set error handler (with retry logic)
    this.websocketConnection.onerror = (e) => {
      if (this.websocketRetriesRemaining > 0){
        --this.websocketRetriesRemaining
        setTimeout(this.openWebsocket.bind(this), 1000)
      } else {
        this.setState({ connectionStatus: STATUS_FAILED })
      }
    }

    // set success handler
    this.websocketConnection.onopen = () => {
      this.websocketRetriesRemaining = WEBSOCKET_RETRY_COUNT
      this.setState({ connectionStatus: STATUS_CONNECTED })
      this.websocketHasOpened = true
    }

    // onclose fires on conn close and failed open
    this.websocketConnection.onclose = () => {
      // if connection failed
      if (!this.websocketHasOpened) return

      // if connection was otherwise closed
      this.setState({ connectionStatus: STATUS_DISCONNECTED })
      delete this.websocketConnection
    }

    // set incoming message handler
    this.websocketConnection.onmessage = (msg) => this.dispatchLog(msg.data)
  }

  dispatchLog(data) {
    // temporary parser
    let log = data.toString()
    // debugger
    let [subsys, ...messageComponents] = log.split(":")
    let message = messageComponents.join(' ')

    let logChannel = 'log:' + (SUBSYSTEMS.includes(subsys) ? subsys : 'misc')


    this.logEmitter.emit(logChannel, {
      message: messageComponents,
      logLevel: 0
    })
  }

  render() {
    let subsystemElements = SUBSYSTEMS.map(subsys => {
      return <Subsystem key={subsys} subsys={subsys} logEmitter={this.logEmitter} />
    })
    return (
      <div>
        <ConnectionStatusIndicator status={this.state.connectionStatus} />
        <div className="Subsystems">
          {subsystemElements}
        </div>
      </div>
    )
  }
}

export { WebKitten }