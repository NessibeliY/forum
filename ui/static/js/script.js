// delete button for post
const dropdowns = document.querySelectorAll(".dropdown");
var formDel = document.querySelector(".select_menu")
var options = document.querySelector(".options")

// var list =  document.getElementById("list")


// Скрытие/отображение dropdown при клике
dropdowns.forEach(dropdown => {
  dropdown.addEventListener('click', (event) => {
      event.stopPropagation(); // Остановка всплытия события
      const targetId = event.target.getAttribute('data-target');
      if (targetId) {
          const formDel = document.getElementById(targetId);
          if (formDel) {
              // Закрываем все dropdown перед открытием нового
              document.querySelectorAll('.select_menu').forEach(menu => {
                  if (menu !== formDel) {
                      menu.style.display = 'none';
                  }
              });

              // Переключаем состояние текущего dropdown
              formDel.style.display = formDel.style.display === 'flex' ? 'none' : 'flex';
          }
      }
  });
});

// Закрытие dropdown при клике вне его
document.addEventListener('click', (event) => {
  const forms = document.querySelectorAll('.select_menu');
  forms.forEach(form => {
      if (!form.contains(event.target)) {
          form.style.display = 'none';
      }
  });
});


// Добавляем обработчик на кнопку "Delete"
document.querySelectorAll('.del_btn_modal').forEach(button => {
  button.addEventListener('click', (event) => {
      event.stopPropagation(); // Остановка всплытия события

      // Логика открытия модального окна (замените под ваш функционал)
      const modalId = button.getAttribute('data-modal');
      if (modalId) {
          const modal = document.getElementById(modalId);
          if (modal) {
              modal.style.display = 'block'; // Показываем модальное окно
          }
      }

      // Закрываем все dropdown
      document.querySelectorAll('.select_menu').forEach(menu => {
          menu.style.display = 'none';
      });
  });
});

// Добавляем обработчик на кнопку "Update"
document.querySelectorAll('.update_btn_modal').forEach(button => {
  button.addEventListener('click', (event) => {
      event.stopPropagation(); // Остановка всплытия события

      // Логика открытия модального окна (замените под ваш функционал)
      const modalId = button.getAttribute('data-modal');
      if (modalId) {
          const modal = document.getElementById(modalId);
          if (modal) {
              modal.style.display = 'block'; // Показываем модальное окно
          }
      }

      // Закрываем все dropdown
      document.querySelectorAll('.select_menu').forEach(menu => {
          menu.style.display = 'none';
      });
  });
});


var del_btn_modals = document.querySelectorAll(".del_btn_modal");
var close_btn_modals = document.querySelectorAll("[data-close]");
var del_Modal_Btns = document.querySelectorAll(".del_Modal_Btn");
var update_btn_modals = document.querySelectorAll(".update_btn_modal")

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

update_btn_modals.forEach((btn) => {
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
})

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

