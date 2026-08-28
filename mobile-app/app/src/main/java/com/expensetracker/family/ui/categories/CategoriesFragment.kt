package com.expensetracker.family.ui.categories

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.expensetracker.family.databinding.FragmentCategoriesBinding
import com.expensetracker.family.network.ApiClient
import com.expensetracker.family.network.ExpenseApi
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch

class CategoriesFragment : Fragment() {

    private lateinit var binding: FragmentCategoriesBinding
    private var categories: List<Map<String, Any>> = emptyList()

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentCategoriesBinding.inflate(inflater, container, false)
        setupRecyclerView()
        fetchCategories()
        return binding.root
    }

    private fun setupRecyclerView() {
        binding.rvCategories.layoutManager = LinearLayoutManager(requireContext())
        binding.rvCategories.adapter = CategoriesAdapter(categories) { categoryName ->
            showCategoryTransactions(categoryName)
        }
    }

    private fun fetchCategories() {
        lifecycleScope.launch(Dispatchers.IO) {
            withContext(Dispatchers.Main) {
                val api = ApiClient.getApiService()
                api.listCategories().enqueue(object : retrofit2.Callback<Map<String, Any>> {
                    override fun onResponse(
                        call: retrofit2.Call<Map<String, Any>>,
                        response: retrofit2.Response<Map<String, Any>>
                    ) {
                        if (response.isSuccessful && response.body() != null) {
                            categories = response.body()!!["data"] as? List<Map<String, Any>> ?: emptyList()
                            binding.rvCategories.adapter = CategoriesAdapter(categories) { categoryName ->
                                showCategoryTransactions(categoryName)
                            }
                        }
                    }

                    override fun onFailure(
                        call: retrofit2.Call<Map<String, Any>>,
                        t: Throwable
                    ) {
                        android.widget.Toast.makeText(requireContext(),
                            "Failed to load categories", android.widget.Toast.LENGTH_SHORT).show()
                    }
                })
            }
        }
    }

    private fun showCategoryTransactions(categoryName: String) {
        android.widget.Toast.makeText(
            requireContext(),
            "Transactions for: $categoryName",
            android.widget.Toast.LENGTH_SHORT).show()
        // TODO: Fetch transactions filtered by this category
    }
}