import React, {useEffect} from 'react'
import {useDispatch, useSelector} from 'react-redux'
import {isReadyForConfiguration} from 'selectors/configuration'
import {resetConfiguration} from 'reducers/configuration'
import InfoView from 'components/infoView'
import ConfigurationForm from 'components/configurationForm'
import ConfigurationView from 'components/configurationView'

import './main-generator.css'

const MainGenerator = () => {
  const dispatch = useDispatch()
  const readyForConfig = useSelector(isReadyForConfiguration)

  useEffect(() => {
    return () => dispatch(resetConfiguration())
  }, [dispatch])

  return (
    <div className="main-generator">
      <div className="main-generator-form-wrapper">
        <h4 className="main-generator-form-subtitle">
          请填写您的系统参数
        </h4>
        <ConfigurationForm />
      </div>
      <div className="main-generator-result-wrapper">
        {readyForConfig ? <ConfigurationView /> : <InfoView />}
      </div>
    </div>
  )
}

export default MainGenerator
