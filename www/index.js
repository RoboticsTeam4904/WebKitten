const React = require('react')
const ReactDOM = require('react-dom')

import { WebKitten } from './components/WebKitten'

ReactDOM.render(
  <div>
    <h1>WebKitten</h1>
    <WebKitten />
  </div>,
  document.getElementsByClassName('webkitten-container')[0]
);