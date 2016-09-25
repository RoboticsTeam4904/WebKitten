const React = require('react')
const ReactDOM = require('react-dom')

import { WebKitten } from './components/WebKitten'

ReactDOM.render(
  <div>
    <WebKitten />
  </div>,
  document.getElementsByClassName('Webkitten-container')[0]
);