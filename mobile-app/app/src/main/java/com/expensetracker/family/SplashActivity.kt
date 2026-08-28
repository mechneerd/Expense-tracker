package com.expensetracker.family

import android.content.Intent
import android.os.Bundle
import android.os.Handler
import androidx.appcompat.app.AppCompatActivity
import com.expensetracker.family.databinding.ActivitySplashBinding

class SplashActivity : AppCompatActivity() {

    private lateinit var binding: ActivitySplashBinding
    private val SPLASH_DURATION = 3000

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivitySplashBinding.inflate(layoutInflater)
        setContentView(binding.root)

        // Check if user already has valid JWT token
        val userToken = android.preference.Preferences.getDefaultSharedPreferences(this)
            .getString("jwt_token", "")

        if (userToken.isNotEmpty()) {
            // User is logged in, go to home
            navigateToHome()
        } else {
            // Show splash for duration, then go to login
            Handler().postDelayed(Runnable {
                navigateToLogin()
            }, SPLASH_DURATION)
        }
    }

    private fun navigateToHome() {
        val intent = Intent(this, MainActivity::class.java)
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)
        startActivity(intent)
        finish()
    }

    private fun navigateToLogin() {
        val intent = Intent(this, LoginOtpActivity::class.java)
        startActivity(intent)
        finish()
    }
}