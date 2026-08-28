package com.expensetracker.family

import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.*
import java.util.concurrent.TimeUnit

interface ApiService {

    @FormUrlEncoded
    @POST("api/v1/auth/google")
    fun googleLogin(
        @Field("email") email: String,
        @Field("first_name") firstName: String,
        @Field("last_name") lastName: String
    ): retrofit2.Call<Map<String, Any>>

    @FormUrlEncoded
    @POST("api/v1/auth/verify-otp")
    fun verifyOtp(
        @Field("email") email: String,
        @Field("otp") otp: String
    ): retrofit2.Call<Map<String, Any>>

    @FormUrlEncoded
    @POST("api/v1/auth/refresh")
    fun refreshToken(
        @Field("email") email: String,
        @Field("refresh_token") refreshToken: String
    ): retrofit2.Call<Map<String, Any>>

    @GET("api/v1/users/me")
    fun getProfile(
        @Header("Authorization") token: String
    ): retrofit2.Call<Map<String, Any>>

    @GET("api/v1/families/me")
    fun getFamilyMembers(
        @Header("Authorization") token: String,
        @Query("family_id") familyId: String
    ): retrofit2.Call<Map<String, Any>>

    @GET("api/v1/transactions/me")
    fun getMyTransactions(
        @Header("Authorization") token: String
    ): retrofit2.Call<Map<String, Any>>

    @GET("api/v1/categories")
    fun getCategories(
        @Header("Authorization") token: String,
        @Query("type") type: String
    ): retrofit2.Call<Map<String, Any>>

    @GET("api/v1/payment-methods")
    fun getPaymentMethods(
        @Header("Authorization") token: String
    ): retrofit2.Call<Map<String, Any>>

    @GET("api/v1/upi-apps")
    fun getUpiApps(
        @Header("Authorization") token: String
    ): retrofit2.Call<Map<String, Any>>
}

object ApiClient {

    // Change this to your backend URL
    // For emulator: http://10.0.2.2:8080
    // For device: http://YOUR_COMPUTER_IP:8080
    // For production: https://your-railway-app.up.railway.app
    private const val BASE_URL = "http://10.0.2.2:8080"

    private val loggingInterceptor = HttpLoggingInterceptor().apply {
        level = HttpLoggingInterceptor.Level.BODY
    }

    private val httpClient = OkHttpClient.Builder()
        .addInterceptor(loggingInterceptor)
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .build()

    private val retrofit = Retrofit.Builder()
        .baseUrl(BASE_URL)
        .client(httpClient)
        .addConverterFactory(GsonConverterFactory.create())
        .build()

    fun getApiService(): ApiService {
        return retrofit.create(ApiService::class.java)
    }
}
