package com.expensetracker.family.ui.dashboard

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.expensetracker.family.databinding.FragmentDashboardBinding
import com.expensetracker.family.network.ApiClient
import com.expensetracker.family.network.ExpenseApi
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import android.preference.Preferences

class HomeFragment : Fragment() {

    private lateinit var binding: FragmentDashboardBinding
    private var transactions: List<Map<String, Any>> = emptyList()
    private var isFamilyHead: Boolean = false

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentDashboardBinding.inflate(inflater, container, false)
        setupRecyclerView()
        fetchData()
        return binding.root
    }

    private fun setupRecyclerView() {
        binding.rvTransactions.layoutManager = LinearLayoutManager(requireContext())
        binding.rvTransactions.adapter = TransactionAdapter(transactions)
    }

    private fun fetchData() {
        lifecycleScope.launch(Dispatchers.IO) {
            withContext(Dispatchers.Main) {
                val api = ApiClient.getApiService()
                val userPrefs = Preferences.getDefaultSharedPreferences(requireContext())
                val userRole = userPrefs.getString("family_role", "")

                // Check if user is family head
                isFamilyHead = userRole == "HEAD"
                binding.tvFamilyRole.text = when {
                    isFamilyHead -> "Family Head"
                    else -> "Family Member"
                }

                // Show/dashboard only for head
                binding.dashboardContainer.isVisible = isFamilyHead

                // Fetch my transactions
                api.listMyTransactions().enqueue(object : retrofit2.Callback<Map<String, Any>> {
                    override fun onResponse(
                        call: retrofit2.Call<Map<String, Any>>,
                        response: retrofit2.Response<Map<String, Any>>
                    ) {
                        if (response.isSuccessful && response.body() != null) {
                            transactions = response.body()!!["data"] as? List<Map<String, Any>> ?: emptyList()
                            binding.rvTransactions.adapter = TransactionAdapter(transactions)
                        }
                    }

                    override fun onFailure(
                        call: retrofit2.Call<Map<String, Any>>,
                        t: Throwable
                    ) {
                        android.widget.Toast.makeText(requireContext(),
                            "Failed to load transactions", android.widget.Toast.LENGTH_SHORT).show()
                    }
                })

                // Fetch dashboard data (head only)
                if (isFamilyHead) {
                    api.getDashboard("monthly").enqueue(object : retrofit2.Callback<Map<String, Any>> {
                        override fun onResponse(
                            call: retrofit2.Call<Map<String, Any>>,
                            response: retrofit2.Response<Map<String, Any>>
                        ) {
                            if (response.isSuccessful && response.body() != null) {
                                val dashboard = response.body()!!
                                binding.tvTotalIncome.text =
                                    formatCurrency(dashboard["total_income"] as? Double ?: 0)
                                binding.tvTotalExpense.text =
                                    formatCurrency(dashboard["total_expense"] as? Double ?: 0)
                                val balance =
                                    (dashboard["total_income"] as? Double ?: 0) -
                                    (dashboard["total_expense"] as? Double ?: 0)
                                binding.tvBalance.text = formatCurrency(balance)
                                binding.tvThisMonthIncome.text =
                                    formatCurrency(dashboard["this_month_income"] as? Double ?: 0)
                                binding.tvThisMonthExpense.text =
                                    formatCurrency(dashboard["this_month_expense"] as? Double ?: 0)
                            }
                        }

                        override fun onFailure(
                            call: retrofit2.Call<Map<String, Any>>,
                            t: Throwable
                        ) {
                            // Dashboard errors non-fatal for now
                        }
                    })
                }
            }
        }
    }

    private fun formatCurrency(amount: Double): String {
        return "${amount.toStringAsFixed(2)} ₹"
    }
}