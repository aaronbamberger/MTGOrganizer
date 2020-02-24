import React from 'react';
import './Keyrune-master/css/keyrune.css'

function SetSymbol(props) {
  return (
    <span
      title={props.setName}
      className={"ss ss-" + props.setCode.toLowerCase()} />
  );
}

export {SetSymbol};
