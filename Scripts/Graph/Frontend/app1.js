document.getElementById("getUpdates").addEventListener("click", function(){
  var cantidad = document.getElementById("Cantidad").value;
  console.log(cantidad);
  if(cantidad>500){
    alert("Debe ser menos de 500 números");
    return
  }
  axios.get(`/generate/${cantidad}`)
})
document.getElementById("comenzar").addEventListener("click", function(){

  axios.get(`/start`)
})


function array_move(arr, old_index, new_index) {
  if (new_index >= arr.length) {
      var k = new_index - arr.length + 1;
      while (k--) {
          arr.push(undefined);
      }
  }
  arr.splice(new_index, 0, arr.splice(old_index, 1)[0]);
};

// returns [2, 1, 3]
console.log(array_move([1, 2, 3], 0, 1)); 

document.addEventListener("DOMContentLoaded", function() {
  console.log("Hola");
  axios.get(`/load`)
});


const pusher = new Pusher('7befe6ab035a03a2ada9', {
    cluster: 'us2',
    encrypted: true
});

const channel = pusher.subscribe('ArrayChannel');


channel.bind("insertion",data =>{
  var a = ChartInsertion.data.datasets[0].data[data[0]];


    var from = data[0];     

    var to = data[1];  
  
    // Store the moved element in a temp  
    // variable 
    var temp = ChartInsertion.data.datasets[0].data[from];  
      
    // shift elements forward  
    var i; 
    for (i = from; i >= to; i--)  
        { 
          ChartInsertion.data.datasets[0].data[i] = ChartInsertion.data.datasets[0].data[i - 1];  
        } 
      
    // Insert moved element at position   
    ChartInsertion.data.datasets[0].data[to] = temp;  

  ChartInsertion.update();
  
})

channel.bind("bubble",data =>{
  var from = data[0]
  var to= data[1]
  var a = ChartBubble.data.datasets[0].data[from];
  ChartBubble.data.datasets[0].data[from] = ChartBubble.data.datasets[0].data[to];
  ChartBubble.data.datasets[0].data[to] = a ;
  ChartBubble.update();
  
})

channel.bind("quick",data =>{
  console.log(data)
  var from = data[0]
  var to= data[1]
  var a = ChartQuick.data.datasets[0].data[from];
  ChartQuick.data.datasets[0].data[from] = ChartQuick.data.datasets[0].data[to];
  ChartQuick.data.datasets[0].data[to] = a ;
  ChartQuick.update();
  
})

channel.bind("a",data =>{
  console.log("llegó");
})

channel.bind('addNumber', data => {


    data.forEach(function (value, i){
        
      ChartBubble.data.labels.push(i+1);
      ChartBubble.data.datasets[0].data.push(value);
      ChartInsertion.data.labels.push(i+1);
      ChartInsertion.data.datasets[0].data.push(value);
      ChartHeap.data.labels.push(i+1);
      ChartHeap.data.datasets[0].data.push(value);
      ChartQuick.data.labels.push(i+1);
      ChartQuick.data.datasets[0].data.push(value);
    });
    ChartBubble.update();
    ChartInsertion.update();
    ChartHeap.update();
    ChartQuick.update();
});

function renderChart() {
  var canvasBubble = document.getElementById("bubble").getContext("2d");
  var canvasInsertion = document.getElementById("insertion").getContext("2d");
  var canvasHeap = document.getElementById("heap").getContext("2d");
  var canvasQuick = document.getElementById("quick").getContext("2d");
  var options = { animation: { duration: 0 }};

  ChartBubble = new Chart(canvasBubble, {
    type: "bar",
    data: getconfigChart(),
    options: options
  });
  ChartInsertion = new Chart(canvasInsertion, {
    type: "bar",
    data: getconfigChart(),
    options: options
  });
  ChartHeap = new Chart(canvasHeap, {
    type: "bar",
    data: getconfigChart(),
    options: options
  });
  ChartQuick = new Chart(canvasQuick, {
    type: "bar",
    data: getconfigChart(),
    options: options
  });

};
function getconfigChart(){
  return {
    labels: [],
    datasets: [
       {
          label: "Magnitud de los números",
          fill: false,
          lineTension: 0.1,
          backgroundColor: "rgba(75,192,192,0.4)",
          borderColor: "rgba(75,192,192,1)",
          borderCapStyle: 'butt',
          borderDash: [],
          borderDashOffset: 0.0,
          borderJoinStyle: 'miter',
          pointBorderColor: "rgba(75,192,192,1)",
          pointBackgroundColor: "#fff",
          pointBorderWidth: 1,
          pointHoverRadius: 5,
          pointHoverBackgroundColor: "rgba(75,192,192,1)",
          pointHoverBorderColor: "rgba(220,220,220,1)",
          pointHoverBorderWidth: 2,
          pointRadius: 1,
          pointHitRadius: 10,
          data: [],
          spanGaps: false,
       }
    ]
  };
}


renderChart()
