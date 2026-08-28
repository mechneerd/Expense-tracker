package com.expensetracker.family

import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.expensetracker.family.databinding.ActivityLoginOtpBinding
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class LoginOtpActivity : AppCompatActivity() {

    private lateinit var binding: ActivityLoginOtpBinding
    private var timerHandler: Handler? = null
    private var remainingTime: Long = 60000
    private val prefs by lazy { getSharedPreferences("expense_prefs", Context.MODE_PRIVATE) }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityLoginOtpBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val savedEmail = prefs.getString("user_email", "") ?: ""
        binding.etEmail.setText(savedEmail)

        binding.btnGoogle.setOnClickListener {
            googleLogin()
        }

        binding.btnSendOtp.setOnClickListener {
            sendOtp()
        }

        binding.tvResend.setOnClickListener {
            sendOtp()
        }

        startResendTimer()
    }

    private fun googleLogin() {
        val email = binding.etEmail.text.toString().trim()
        if (email.isEmpty()) {
            binding.etEmail.error = "Enter email"
            return
        }

        binding.btnGoogle.isEnabled = false

        val api = ApiClient.getApiService()
        api.googleLogin(email, "", "").enqueue(object : Callback<Map<String, Any>> {
            override fun onResponse(
                call: Call<Map<String, Any>>,
                response: Response<Map<String, Any>>
            ) {
                binding.btnGoogle.isEnabled = true

                if (response.isSuccessful) {
                    val body = response.body()
                    val data = body?.get("data") as? Map<*, *>
                    val token = data?.get("token")?.toString() ?: ""

                    prefs.edit()
                        .putString("jwt_token", token)
                        .putString("user_email", email)
                        .apply()

                    Toast.makeText(this@LoginOtpActivity, "OTP sent to $email", Toast.LENGTH_LONG).show()
                } else {
                    Toast.makeText(this@LoginOtpActivity, "Login failed", Toast.LENGTH_LONG).show()
                }
            }

            override fun onFailure(call: Call<Map<String, Any>>, t: Throwable) {
                binding.btnGoogle.isEnabled = true
                Toast.makeText(this@LoginOtpActivity, "Network error: ${t.message}", Toast.LENGTH_LONG).show()
            }
        })
    }

    private fun sendOtp() {
        val email = binding.etEmail.text.toString().trim()
        if (email.isEmpty()) {
            binding.etEmail.error = "Enter email"
            return
        }

        val otp = "123456" // For testing - in real app, user enters OTP from email

        val api = ApiClient.getApiService()
        api.verifyOtp(email, otp).enqueue(object : Callback<Map<String, Any>> {
            override fun onResponse(
                call: Call<Map<String, Any>>,
                response: Response<Map<String, Any>>
            ) {
                if (response.isSuccessful) {
                    val body = response.body()
                    val data = body?.get("data") as? Map<*, *>
                    val token = data?.get("token")?.toString() ?: ""
                    val refreshToken = data?.get("refresh")?.toString() ?: ""

                    prefs.edit()
                        .putString("jwt_token", token)
                        .putString("refresh_token", refreshToken)
                        .putString("user_email", email)
                        .apply()

                    Toast.makeText(this@LoginOtpActivity, "OTP verified!", Toast.LENGTH_SHORT).show()
                    navigateToHome()
                } else {
                    Toast.makeText(this@LoginOtpActivity, "Invalid OTP", Toast.LENGTH_LONG).show()
                }
            }

            override fun onFailure(call: Call<Map<String, Any>>, t: Throwable) {
                Toast.makeText(this@LoginOtpActivity, "Network error", Toast.LENGTH_LONG).show()
            }
        })
    }

    private fun startResendTimer() {
        timerHandler = Handler(Looper.getMainLooper())
        remainingTime = 60000

        val timerRunnable = object : Runnable {
            override fun run() {
                val minutes = (remainingTime / 60000).toString()
                val seconds = ((remainingTime % 60000) / 1000).toString()
                binding.tvResend.text = "Resend OTP ($minutes:$seconds)"

                if (remainingTime > 0) {
                    remainingTime -= 1000
                    timerHandler?.postDelayed(this, 1000)
                    binding.tvResend.visibility = View.VISIBLE
                    binding.btnSendOtp.visibility = View.GONE
                } else {
                    binding.tvResend.visibility = View.GONE
                    binding.btnSendOtp.visibility = View.VISIBLE
                    timerHandler = null
                }
            }
        }

        timerHandler?.postDelayed(timerRunnable, 1000)
    }

    private fun navigateToHome() {
        val intent = Intent(this, MainActivity::class.java)
        intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
        startActivity(intent)
        finish()
    }

    override fun onDestroy() {
        super.onDestroy()
        timerHandler?.removeCallbacksAndMessages(null)
    }
}
