document.addEventListener('DOMContentLoaded', function() {
    const burger = document.getElementById('burgerBtn');
    const mobileMenu = document.getElementById('mobileMenu');
    const mobileLinks = document.querySelectorAll('.mobile-link');

    if (burger) {
        burger.addEventListener('click', function() {
            burger.classList.toggle('active');
            mobileMenu.classList.toggle('active');
        });
    }

    mobileLinks.forEach(function(link) {
        link.addEventListener('click', function() {
            burger.classList.remove('active');
            mobileMenu.classList.remove('active');
        });
    });


    const statNumbers = document.querySelectorAll('.stat-number');

    const animateNumbers = function(entries) {
        entries.forEach(function(entry) {
            if (entry.isIntersecting) {
                const el = entry.target;
                const target = parseInt(el.getAttribute('data-target'));
                let current = 0;
                const increment = target / 50;

                const timer = setInterval(function() {
                    current += increment;
                    if (current >= target) {
                        el.textContent = target;
                        clearInterval(timer);
                    } else {
                        el.textContent = Math.floor(current);
                    }
                }, 30);
            }
        });
    };

    const observer = new IntersectionObserver(animateNumbers, {
        threshold: 0.5
    });

    statNumbers.forEach(function(num) {
        observer.observe(num);
    });


    const featureCards = document.querySelectorAll('.feature-card');

    const showCards = function(entries) {
        entries.forEach(function(entry) {
            if (entry.isIntersecting) {
                const card = entry.target;
                const delay = parseInt(card.getAttribute('data-delay')) || 0;
                setTimeout(function() {
                    card.classList.add('visible');
                }, delay);
            }
        });
    };

    const cardObserver = new IntersectionObserver(showCards, {
        threshold: 0.15
    });

    featureCards.forEach(function(card) {
        cardObserver.observe(card);
    });


    const modal = document.getElementById('modal');
    const modalClose = document.getElementById('modalClose');
    const modalConfirm = document.getElementById('modalConfirm');
    const selectBtns = document.querySelectorAll('.select-plan');
    
    const modalFormState = document.getElementById('modalFormState');
    const modalSuccessState = document.getElementById('modalSuccessState');
    
    const modalSuccessTitle = document.getElementById('modalSuccessTitle');
    const modalSuccessText = document.getElementById('modalSuccessText');
    const modalSuccessSub = document.getElementById('modalSuccessSub');

    const showSuccessModal = function(titleText, bodyText, subText) {
        if (!modal) return;
        
        if (modalSuccessTitle) modalSuccessTitle.innerHTML = titleText;
        if (modalSuccessText) modalSuccessText.textContent = bodyText;
        if (modalSuccessSub) modalSuccessSub.textContent = subText;

        if (modalFormState) modalFormState.style.display = 'none';
        if (modalSuccessState) modalSuccessState.style.display = 'block';

        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    };

    selectBtns.forEach(function(btn) {
        btn.addEventListener('click', function() {
            const plan = btn.getAttribute('data-plan');
            
            document.getElementById('modalPlanName').textContent = plan;
            document.getElementById('modalPlanInput').value = plan;

            modalForm.reset();
            modalName.classList.remove('invalid');
            modalNameError.classList.remove('visible');
            modalPhone.classList.remove('invalid');
            modalPhoneError.classList.remove('visible');
            modalEmail.classList.remove('invalid');
            modalEmailError.classList.remove('visible');

            if (modalFormState) modalFormState.style.display = 'block';
            if (modalSuccessState) modalSuccessState.style.display = 'none';

            modal.classList.add('active');
            document.body.style.overflow = 'hidden';
        });
    });

    const closeModal = function() {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    };

    if (modalClose) modalClose.addEventListener('click', closeModal);
    if (modalConfirm) modalConfirm.addEventListener('click', closeModal);

    modal.addEventListener('click', function(e) {
        if (e.target === modal) {
            closeModal();
        }
    });

    
    const validateName = function(name) {
        return name.trim().length >= 2;
    };

    const validatePhone = function(phone) {
        const phoneRegex = /^(\+7|7|8)?[\s\-]?\(?[0-9]{3}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$/;
        return phoneRegex.test(phone.trim());
    };

    const validateEmail = function(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email.trim());
    };

    const modalForm = document.getElementById('modalForm');
    const modalName = document.getElementById('modalName');
    const modalPhone = document.getElementById('modalPhone');
    const modalEmail = document.getElementById('modalEmail');
    
    const modalNameError = document.getElementById('modalNameError');
    const modalPhoneError = document.getElementById('modalPhoneError');
    const modalEmailError = document.getElementById('modalEmailError');
    const modalSubmitBtn = document.getElementById('modalSubmitBtn');

    if (modalForm) {
        modalForm.addEventListener('submit', function(e) {
            e.preventDefault();

            let isValid = true;

            if (!validateName(modalName.value)) {
                modalName.classList.add('invalid');
                modalNameError.classList.add('visible');
                isValid = false;
            } else {
                modalName.classList.remove('invalid');
                modalNameError.classList.remove('visible');
            }

            if (!validatePhone(modalPhone.value)) {
                modalPhone.classList.add('invalid');
                modalPhoneError.classList.add('visible');
                isValid = false;
            } else {
                modalPhone.classList.remove('invalid');
                modalPhoneError.classList.remove('visible');
            }

            if (!validateEmail(modalEmail.value)) {
                modalEmail.classList.add('invalid');
                modalEmailError.classList.add('visible');
                isValid = false;
            } else {
                modalEmail.classList.remove('invalid');
                modalEmailError.classList.remove('visible');
            }

            if (!isValid) return;

            modalSubmitBtn.disabled = true;
            modalSubmitBtn.textContent = 'Отправка...';

            const formData = {
                plan: document.getElementById('modalPlanInput').value,
                name: modalName.value,
                phone: modalPhone.value,
                email: modalEmail.value
            };

            fetch(`/api/order`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
            .then(function(response) {
                return response.json();
            })
            .then(function(data) {
                modalSubmitBtn.disabled = false;
                modalSubmitBtn.textContent = 'Отправить';

                if (data.success) {
                    modalForm.reset();
                    showSuccessModal('🎉 Спасибо!', 'Заявка на тариф успешно оформлена', 'Менеджер свяжется с вами в течение 15 минут.');
                } else {
                    const errorText = data.message || 'Произошла ошибка на стороне сервера.';
                    showSuccessModal('⚠️ Ошибка', errorText, 'Пожалуйста, проверьте введённые данные.');
                }
            })
            .catch(function(error) {
                modalSubmitBtn.disabled = false;
                modalSubmitBtn.textContent = 'Отправить';
                console.error('Ошибка:', error);
                showSuccessModal('❌ Ошибка подключения', 'Не удалось связаться с сервером.', 'Убедитесь, что сервер Go запущен.');
            });
        });
    }

    const contactForm = document.getElementById('contactForm');
    const formName = document.getElementById('formName');
    const formPhone = document.getElementById('formPhone');
    const formMessage = document.getElementById('formMessage');
    
    const nameError = document.getElementById('nameError');
    const phoneError = document.getElementById('phoneError');
    const formSubmitBtn = document.getElementById('formSubmitBtn');

    if (contactForm) {
        contactForm.addEventListener('submit', function(e) {
            e.preventDefault();

            let isValid = true;

            if (!validateName(formName.value)) {
                formName.classList.add('invalid');
                nameError.classList.add('visible');
                isValid = false;
            } else {
                formName.classList.remove('invalid');
                nameError.classList.remove('visible');
            }

            if (!validatePhone(formPhone.value)) {
                formPhone.classList.add('invalid');
                phoneError.classList.add('visible');
                isValid = false;
            } else {
                formPhone.classList.remove('invalid');
                phoneError.classList.remove('visible');
            }

            if (!isValid) return;

            formSubmitBtn.disabled = true;
            formSubmitBtn.textContent = 'Отправка...';

            const formData = {
                name: formName.value,
                phone: formPhone.value,
                message: formMessage.value
            };

            fetch(`/api/contact`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
            .then(function(response) {
                return response.json();
            })
            .then(function(data) {
                formSubmitBtn.disabled = false;
                formSubmitBtn.textContent = 'Отправить';

                if (data.success) {
                    contactForm.reset();
                    showSuccessModal('🎉 Спасибо!', 'Ваша заявка успешно отправлена', 'Менеджер свяжется с вами в течение 15 минут.');
                } else {
                    const errorText = data.message || 'Произошла ошибка на стороне сервера.';
                    showSuccessModal('⚠️ Ошибка', errorText, 'Пожалуйста, проверьте введённые данные.');
                }
            })
            .catch(function(error) {
                formSubmitBtn.disabled = false;
                formSubmitBtn.textContent = 'Отправить';
                console.error('Ошибка:', error);
                showSuccessModal('❌ Ошибка подключения', 'Не удалось связаться с сервером.', 'Убедитесь, что сервер Go запущен.');
            });
        });
    }

});