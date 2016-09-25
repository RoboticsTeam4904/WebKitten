const React = require('react')
const Immutable = require('immutable')
import { LogBox } from './index'
import { LEVEL_DEBUG } from './LogBox'
const lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

class Subsystem extends React.Component {
	constructor(props) {
		super(props)

		this.state = {
			logMessages: Immutable.List()
		}
	}

	componentWillMount() {
		let emitter = this.props.logEmitter
		this.logSubscription = emitter.addListener(`log:${this.subsystem()}`, this.logMessageHandler.bind(this))
	}

	componentWillUnmount() {
		this.logSubscription.remove()
	}

	logMessageHandler(logMessage) {
		this.setState({
			logMessages: this.state.logMessages.push(logMessage)
		})
	}

  render() {
    return (
    	<div className="Subsystem-container">
    		<h3 className="Subsystem-label">Subsystem: {this.subsystem()}</h3>
  			<LogBox logs={this.state.logMessages} highestLogLevelToDisplay={LEVEL_DEBUG} />
  		</div>
		)
  }

  subsystem() {
    return this.props.subsys
  }
}

export { Subsystem }