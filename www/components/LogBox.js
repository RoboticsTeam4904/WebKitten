const React = require('react')
const Immutable = require('immutable')

const LEVEL_WTF = -1
const LEVEL_FATAL = 0
const LEVEL_ERROR = 1
const LEVEL_WARN = 2
const LEVEL_INFO = 3
const LEVEL_DEBUG = 4
const LOG_LEVELS = {
  [LEVEL_WTF]: 'WTF',
  [LEVEL_FATAL]: 'FATAL',
  [LEVEL_ERROR]: 'ERROR',
  [LEVEL_WARN]: 'WARN',
  [LEVEL_INFO]: 'INFO',
  [LEVEL_DEBUG]: 'DEBUG'
}

class LogBox extends React.Component {

  componentDidUpdate() {
    this.scrollToLatestLog()
  }

  render() {
    return (
      <div ref="logBox" className="Subsystem-log-box">
        {this.visibleLogElements()}
      </div>
    )
  }

  visibleLogMessages() {
    // return this.props.logs.filter(l => l.level <= this.props.highestLogLevelToDisplay)
    return this.props.logs
  }

  scrollToLatestLog() {
    let logBox = this.refs.logBox
    logBox.scrollTop = logBox.scrollHeight
  }

  visibleLogElements() {
    return this.visibleLogMessages().map((message, index) => {
      return (
        <div className="Subsystem-log-item" key={index}>{message}</div>
      )
    })
  }
}

export { LogBox }