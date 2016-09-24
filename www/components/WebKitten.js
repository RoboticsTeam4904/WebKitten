const React = require('react')
import {Subsystem} from './Subsystem'

const SUBSYSTEMS = ['chassis', 'intake', 'flywheel']

class WebKitten extends React.Component {
  render() {
    let subsystemElements = SUBSYSTEMS.map(subsys => {
      return <Subsystem key={subsys} subsys={subsys} />
    })
    return (
      <div>
        {subsystemElements}
      </div>
    )
  }
}

export { WebKitten }