import React, {useState, useEffect} from 'react';
import './App.css';

const { SensorRequest, SensorResponse  } = require("./sensorpb/sensor_pb")
const { SensorClient  } = require("./sensorpb/sensor_grpc_web_pb")

var client = new SensorClient('http://localhost:8000')

export const App = () => {
  const [temp, setTemp] = useState(-9999);
  const [humidity , setHumidity] = useState(-99999)

  const getTemp = () => {
    console.log("called")

    var sensorRequest = new SensorRequest()
    var stream = client.tempSensor(sensorRequest,{})

    stream.on('data', (response: any) =>{
        setTemp(response.getValue())
    });
  };

  const getHumidity = () => {
    var sensorRequest = new SensorRequest()
    var stream = client.humiditySensor(sensorRequest,{})

    stream.on('data',(response: any) => {
      setHumidity(response.getValue())
    })
  }

  useEffect(()=>{
    getTemp()
    getHumidity()
  },[]);

  return (
    <div>
      Temperature : {temp} F
      <br/>
      Humidity : {humidity} %
    </div>
  );
}

export default App;
