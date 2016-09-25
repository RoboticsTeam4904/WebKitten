const React = require('react')

const STATUS_DISCONNECTED = 0
const STATUS_CONNECTING = 1
const STATUS_CONNECTED = 2
const STATUS_FAILED = 3
const STATUS_NAMES = {
  [STATUS_DISCONNECTED]: 'disconnected',
  [STATUS_CONNECTING]: 'connecting',
  [STATUS_CONNECTED]: 'connected',
  [STATUS_FAILED]: 'failed to connect'
}


class ConnectionStatusIndicator extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      connectionStatus: this.props.status
    }
  }

  componentWillReceiveProps(newProps) {
    this.setState({
      connectionStatus: newProps.status
    })
  }

  render() {
    // <span className="connection-status-label">Connection:</span>&nbsp;
    return (
      <div className="ConnectionStatusIndicator">
        <div className={this.statusClassName()}>{this.status()}</div>
      </div>
    )
  }

  status() {
    return STATUS_NAMES[this.state.connectionStatus]
  }

  statusClassName() {
    return `connection-status-value status-${this.status()}`
  }
}
export {
  ConnectionStatusIndicator,
  STATUS_CONNECTED,
  STATUS_CONNECTING,
  STATUS_DISCONNECTED,
  STATUS_FAILED
}