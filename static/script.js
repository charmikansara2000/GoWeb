// select all the buttons
const buttons = document.querySelectorAll('button');
// select the <input type="text" class="display" disabled> element
const display = document.querySelector('.display');

// add eventListener to each button
buttons.forEach(function(button) {
  button.addEventListener('click', calculate);
});

// calculate function
function calculate(event) {
  // current clicked buttons value
  const clickedButtonValue = event.target.value;

  if (clickedButtonValue === '=') {
    // check if the display is not empty then only do the calculation
    if (display.value !== '') {
      console.log(display.value)
      let user = {
				exp: display.value
			};

			async function fet(){
				console.log(user)
        console.log("hiii")
				let response = await fetch('http://localhost:9090/fetch', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json;charset=utf-8'
				},
				body: JSON.stringify(user)
				})
        .then(response => response.json())
        .then(data => 
          //console.log(data)
          document.getElementById("display1").value = data.result
          );

        // const data = response.json();
        // return data.result;

				//let result = await response.json();
				// alert(result.message);
			}
			fet();
      // calculate and show the answer to display
     // display.value = eval(display.value);
    } }else if (clickedButtonValue === 'C') {
    // clear everything on display
    display.value = '';
  } else {
    // otherwise concatenate it to the display
    display.value += clickedButtonValue;
  }
}