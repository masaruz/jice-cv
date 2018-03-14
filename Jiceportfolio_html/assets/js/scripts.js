$(function() {

  "use strict";

  /*===============================================
    Preloader
  ===============================================*/
  $(window).load(function () {
    $("body").addClass("loaded");
  });

  /*===============================================
    Icon Menu
  ===============================================*/
  var windowWidth = $(window).width();
  var iconMenu = $(".icon-menu");
  var menuPosition = iconMenu.offset();

  // Stick Menu after scroll
  if (windowWidth >= 992) {
    $(window).on("scroll", function () {
      if ($(window).scrollTop() >= menuPosition.top) {
        iconMenu.addClass("fixed");
      } else {
        iconMenu.removeClass("fixed");
      }
     });
  }

  /*===============================================
    Sticky Social Media buttons
  ===============================================*/
  var stickySocial = $(".sticky-social");
  var socialMediaPosition = stickySocial.offset();

  if (windowWidth >= 992) {
    $(window).on("scroll", function () {
      if ($(window).scrollTop() >= socialMediaPosition.top) {
        stickySocial.addClass("sticky-social-fixed");
      } else {
        stickySocial.removeClass("sticky-social-fixed");
      }
    });
  }

  /*===============================================
    Toggle Menu
  ===============================================*/
  var toggleBtn = $(".toggle-btn");

  toggleBtn.on("click", function(e) {
    if (iconMenu.hasClass("icon-menu-show")) {
      iconMenu.removeClass("icon-menu-show");
    }
    else {
      iconMenu.addClass("icon-menu-show");
    }
    if (stickySocial.hasClass("sticky-social-show")) {
      stickySocial.removeClass("sticky-social-show");
    }
    else {
      stickySocial.addClass("sticky-social-show");
    }
    e.stopPropagation();
  });

  // Navicon transform into X //
  toggleBtn.on("click", function() {
    if (toggleBtn.hasClass("toggle-close")) {
      toggleBtn.removeClass("toggle-close");
    }
    else {
      toggleBtn.addClass("toggle-close");
    }
  });

  /*===============================================
    Smooth Scrolling
  ===============================================*/
  var htmlBody = $("html,body");
  var smoothLinks = $(".icon-menu ul li a, .homeWrapper .btn-style");

  smoothLinks.on("click", function(e) {
      htmlBody.animate({scrollTop: $(this.hash).offset().top}, 700, "easeInOutQuart");  
    e.preventDefault();
  });

  /*===============================================
    Scroll Spy
  ===============================================*/
  $('body').scrollspy({ 
    target: '.icon-menu'
  });

  /*===============================================
    Typed js
  ===============================================*/
  $(window).load(function () {
    var typed = new Typed("#typed", {
      stringsElement: "#typed-strings",
      startDelay: 100,
      typeSpeed: 30,
      backDelay: 1000,
      backSpeed: 20,
      loop: true,
      loopCount: 2
    });
  });

  /*===============================================
    MixItUp
  ===============================================*/
  $('#mix-container').mixItUp();

  /*===============================================
    Magnific Popup
  ===============================================*/
  $('.lightbox').magnificPopup({ 
    type:'inline',
    fixedContentPos: false,
    removalDelay: 100,
    closeBtnInside: true,
    preloader: false,
    mainClass: 'mfp-fade'
  });
  
  /*===============================================
    Counter
  ===============================================*/
  $(".facts-background").appear(function() {

    var factsCounter = $('.counter');

    factsCounter.each(function () {
      $(this).prop('Counter',0).animate({
          Counter: $(this).text()
      }, {
          duration: 3000,
          easing: 'swing',
          step: function (now) {
              $(this).text(Math.ceil(now));
          }
      });
    });
  },{accX: 0, accY: -10});

  /*===============================================
    Owl Carousel Sliders
  ===============================================*/
  // ===== Clients Slider =====
  $("#clientsSlider").owlCarousel({
    items:3,
    dots:false,
    rewind:true,
    margin:30,
    autoplay:true,
    autoplayHoverPause:true,
    autoplayTimeout:3000, // 3 seconds
    autoplaySpeed:300, // 0.3 seconds
    responsive : {
      // breakpoint from 0 up
      0 : {
        items: 1
      },
      // breakpoint from 768 up
      768 : {
        items: 3
      }
    }
  });
  // Custom Navigation of Clients
  var clientsNavigation = $("#clientsSlider");
  var clientsNext = $("#clientsNext");
  var clientsPrev = $("#clientsPrev");
  // Events
  clientsNext.on("click", function(){
    clientsNavigation.trigger('next.owl.carousel', [300]);
  });
  clientsPrev.on("click", function(){
    clientsNavigation.trigger('prev.owl.carousel', [300]);
  });
  // end Custom Navigation of Blog

  // ===== Blog Slider =====
  $("#blogSlider").owlCarousel({
    items:2,
    dots:false,
    rewind:true,
    responsive : {
      // breakpoint from 0 up
      0 : {
        items: 1
      },
      // breakpoint from 768 up
      768 : {
        items: 2
      }
    }
  });

  // Custom Navigation of Blog
  var blogNavigation = $("#blogSlider");
  var blogNext = $("#next");
  var blogPrev = $("#prev");
  // Events
  blogNext.on("click", function(){
    blogNavigation.trigger('next.owl.carousel', [300]);
  });
  blogPrev.on("click", function(){
    blogNavigation.trigger('prev.owl.carousel', [300]);
  });
  // end Custom Navigation of Blog

  // ===== Testimonial Slider =====
  $("#testimonialSlider").owlCarousel({
    items:1,
    rewind:true,
    margin:30,
    dots:true,
    dotsSpeed:300,
    autoplay:true,
    autoplayHoverPause:true,
    autoplayTimeout:4000, // 4 seconds
    autoplaySpeed:300 // 0.3 seconds
  });

  /*===============================================
    Contact Form
  ===============================================*/
  $("#contactform").on('submit',function(e) {
    var name = $("#name").val();
    var email = $("#email").val();
    var message = $("#message").val();
    var nameId = $("#name");
    var emailId = $("#email");
    var messageId = $("#message");
    if (name === '') {
      nameId.css('border-color','rgba(255, 0, 0, 0.5)');
    }
    if (email === '') {
      emailId.css('border-color','rgba(255, 0, 0, 0.5)');
    }
    if (message === '') {
      messageId.css('border-color','rgba(255, 0, 0, 0.5)');
    }
    else {
      $.ajax({
        url:'contact_form.php',
        data:$(this).serialize(),
        type:'POST',
        success:function(data){
          $("#success").show().fadeIn(1000); //=== Show Success Message==
          $('#contactform').each(function(){
            this.reset();
          });
        },
        error:function(data){
          $("#error").show().fadeIn(1000); //===Show Error Message====
        }
      });
    }
    e.preventDefault();
  });

});