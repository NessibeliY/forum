// delete button for post
var dropdown = document.querySelectorAll(".dropdown")
var formDel = document.querySelector(".select_menu")
var options = document.querySelector(".options")

// var list =  document.getElementById("list")


dropdown.forEach(dropdown => {
  dropdown.addEventListener('click', (event) => {
      const targetId = event.target.getAttribute('data-target');
      if (targetId) {
          const formDel = document.getElementById(targetId);        
          if (formDel) {
              formDel.style.display = formDel.style.display === 'flex' ? 'none' : 'flex';
              console.log("open")
          }
      }
  });
});


document.addEventListener('click', (event) => {
  const forms = document.querySelectorAll('.select_menu');
  forms.forEach(form => {
      if (!form.contains(event.target) && !Array.from(dropdown).some(dropdown => dropdown.contains(event.target))) {
          form.style.display = 'none';
      }
  });
});



var del_btn_modals = document.querySelectorAll(".del_btn_modal");
var close_btn_modals = document.querySelectorAll("[data-close]");
var del_Modal_Btns = document.querySelectorAll(".del_Modal_Btn");


// Обработка нажатия на кнопку "Delete"
del_btn_modals.forEach((btn) => {
  btn.addEventListener("click", (event) => {
    const modalId = btn.getAttribute('data-modal');
    const modal = document.getElementById(modalId);
    if (modal) {
      console.log("this")
      modal.style.display = "flex";
    } else {
      console.error(`No modal found for button with modal ID: ${modalId}`);
    }
  });
});

// Закрытие модального окна при нажатии на крестик или кнопку "no"
close_btn_modals.forEach((btn) => {
  btn.addEventListener("click", () => {
    const modalId = btn.getAttribute('data-close');
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.style.display = "none";
    } else {
      console.error(`No modal found for close button with modal ID: ${modalId}`);
    }
  });
});

del_Modal_Btns.forEach((btn, index) => {
  btn.addEventListener("click", () => {
    if (del_modals[index]) {
      del_modals[index].style.display = "none";
    } else {
      console.error(`No modal found for "no" button at index ${index}`);
    }
  });
});

window.onclick = function(event) {
  if (event.target.classList.contains('modal')) {
    event.target.style.display = "none";
  }
};



function toggleMenu() {
  const menu = document.querySelector('.mobile_menu');
  menu.classList.toggle('active');
}
