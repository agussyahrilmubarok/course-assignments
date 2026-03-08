package com.dicoding.storyapp.features.login

import android.animation.AnimatorSet
import android.animation.ObjectAnimator
import android.content.Intent
import android.os.Build
import android.os.Bundle
import android.os.Handler
import android.view.View
import android.view.WindowInsets
import android.view.WindowManager
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import com.dicoding.storyapp.common.InputValidator
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.common.showDialog
import com.dicoding.storyapp.data.model.LoginRequest
import com.dicoding.storyapp.data.model.UserModel
import com.dicoding.storyapp.databinding.ActivityLoginBinding
import com.dicoding.storyapp.features.ViewModelFactory
import com.dicoding.storyapp.features.home.HomeActivity

class LoginActivity : AppCompatActivity() {

    private lateinit var binding: ActivityLoginBinding

    private val viewModel by viewModels<LoginViewModel> {
        ViewModelFactory.getInstance(applicationContext)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityLoginBinding.inflate(layoutInflater)
        setContentView(binding.root)
        setupView()
        setupAction()
        playAnimation()
    }

    private fun setupView() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
            window.insetsController?.hide(WindowInsets.Type.statusBars())
        } else {
            window.setFlags(
                WindowManager.LayoutParams.FLAG_FULLSCREEN,
                WindowManager.LayoutParams.FLAG_FULLSCREEN
            )
        }
        supportActionBar?.hide()
    }

    private fun setupAction() {
        binding.loginButton.setOnClickListener {
            val email = binding.edLoginEmail.text.toString()
            val password = binding.edLoginPassword.text.toString()

            if (validateInput(email)) {
                doLogin(LoginRequest(email, password))
            }
        }
    }

    private fun playAnimation() {
        ObjectAnimator.ofFloat(binding.imageView, View.TRANSLATION_X, -30f, 30f).apply {
            duration = 6000
            repeatCount = ObjectAnimator.INFINITE
            repeatMode = ObjectAnimator.REVERSE
        }.start()

        val title = ObjectAnimator.ofFloat(binding.titleTextView, View.ALPHA, 1f).setDuration(100)
        val message =
            ObjectAnimator.ofFloat(binding.messageTextView, View.ALPHA, 1f).setDuration(100)
        val emailTextView =
            ObjectAnimator.ofFloat(binding.emailTextView, View.ALPHA, 1f).setDuration(100)
        val emailEditTextLayout =
            ObjectAnimator.ofFloat(binding.emailEditTextLayout, View.ALPHA, 1f).setDuration(100)
        val passwordTextView =
            ObjectAnimator.ofFloat(binding.passwordTextView, View.ALPHA, 1f).setDuration(100)
        val passwordEditTextLayout =
            ObjectAnimator.ofFloat(binding.passwordEditTextLayout, View.ALPHA, 1f).setDuration(100)
        val login = ObjectAnimator.ofFloat(binding.loginButton, View.ALPHA, 1f).setDuration(100)

        AnimatorSet().apply {
            playSequentially(
                title,
                message,
                emailTextView,
                emailEditTextLayout,
                passwordTextView,
                passwordEditTextLayout,
                login
            )
            startDelay = 100
        }.start()
    }

    private fun doLogin(param: LoginRequest) {
        viewModel.login(param).observe(this) {
            when (it) {
                is ResultState.Loading -> {
                    binding.loginButton.isEnabled = false
                }

                is ResultState.Error -> {
                    binding.loginButton.isEnabled = false
                    showDialog(this, "Gagal Masuk", it.message)
                }

                is ResultState.Success -> {
                    binding.loginButton.isEnabled = true
                    showDialog(
                        this,
                        "Masuk Berhasil",
                        it.data.message
                    )
                    val loginResult = it.data.loginResult
                    val currentUser = UserModel(
                        userId = loginResult.userId,
                        name = loginResult.name,
                        email = param.email,
                        password = param.password,
                        token = loginResult.token,
                        isLogged = true
                    )
                    viewModel.savedUser(currentUser)
                    Handler().postDelayed({
                        val moveIntent = Intent(this, HomeActivity::class.java)
                        startActivity(moveIntent)
                        finish()
                    }, 2000)
                }
            }
        }
    }

    private fun validateInput(email: String): Boolean {
        var isValid = true

        if (!InputValidator.isValidEmail(email)) {
            binding.edLoginEmail.error = "Invalig email address"
            isValid = false
        } else {
            binding.edLoginEmail.error = null
        }

        return isValid
    }
}