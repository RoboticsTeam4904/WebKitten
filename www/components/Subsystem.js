const React = require('react')

class Subsystem extends React.Component {
  render() {
    return <h3>Subsystem: {this.subsystem()}</h3>
  }

  subsystem() {
    debugger
    return this.props.subsys
  }
}

export { Subsystem }