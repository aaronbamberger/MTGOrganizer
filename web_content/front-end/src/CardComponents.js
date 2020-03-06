import React from 'react';

import './Keyrune-master/css/keyrune.css';
import {BACKEND_HOSTNAME} from './Constants.js';

function SetSymbol(props) {
  return (
    <span
      title={props.setName}
      className={"ss ss-" + props.setCode.toLowerCase()} />
  );
}

function CardImage(props) {
  const imageStyle = {
    width: props.sizePercent + '%',
    height: props.sizePercent + '%',
    opacity: props.isDisabled ? '50%' : '100%'}


  return (
    <img
      src={'http://' + BACKEND_HOSTNAME + '/static_content/card_face_images/' + props.uuid + '.png'}
      alt={props.name}
      style={imageStyle} />
  );
}

export {SetSymbol, CardImage};
