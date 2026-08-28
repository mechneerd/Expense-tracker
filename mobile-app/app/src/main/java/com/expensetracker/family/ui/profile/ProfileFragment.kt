package com.expensetracker.family.ui.profile

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.expensetracker.family.databinding.FragmentProfileBinding
import com.expensetracker.family.network.ApiClient
import com.expensetracker.family.network.ExpenseApi
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import android.preference.Preferences

class ProfileFragment : Fragment() {

    private lateinit var binding: FragmentProfileBinding
    private var userTransactions: List<Map<String, Any>> = emptyList()

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentProfileBinding.inflate(inflater, container, false)
        setupRecyclerView()
        fetchUserData()
        return binding.root
    }

    private fun setupRecyclerView() {
        binding.rvTransactions.layoutManager = LinearLayoutManager(requireContext())
        binding.rvTransactions.adapter = ProfileTransactionAdapter(userTransactions)
    }

    private fun fetchUserData() {
        lifecycleScope.launch(Dispatchers.IO) {
            withContext(Dispatchers.Main) {
                val api = ApiClient.getApiService()
                val userPrefs = android.preference.Preferences
                    .getDefaultSharedPreferences(requireContext())
                val userRole = userPrefs.getString("family_role", "")

                // Get user profile info from shared preferences
                binding.tvName.text = userPrefs.getString("user_name", "User")
                binding.tvEmail.text = userPrefs.getString("user_email", "email@example.com")
                binding.tvRole.text = when {
                    userRole == "HEAD" -> "Family Head"
                    else -> "Family Member"
                }

                // Get user's own transactions
                api.listMyTransactions().enqueue(object : retrofit2.Callback<Map<String, Any>> {
                    override fun onResponse(
                        call: retrofit2.Call<Map<String, Any>>,
                        response: retrofit2.Response<Map<String, Any>>
                    ) {
                        if (response.isSuccessful && response.body() != null) {
                            userTransactions = response.body()!!["data"] as? List<Map<String, Any>> ?: emptyList()
                            binding.rvTransactions.adapter = ProfileTransactionAdapter(userTransactions)
                        }
                    }

                    override fun onFailure(
                        call: retrofit2.Call<Map<String, Any>>,
                        t: Throwable
                    ) {
                        android.widget.Toast.makeText(requireContext(),
                            "Failed to load user data", android.widget.Toast.LENGTH_SHORT).show()
                    }
                })
            }
        }
    }
}