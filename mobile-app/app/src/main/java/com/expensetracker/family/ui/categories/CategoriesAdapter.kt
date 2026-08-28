package com.expensetracker.family.ui.categories

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.expensetracker.family.R

class CategoriesAdapter(
    private var categories: List<Map<String, Any>>,
    onCategoryClick: (String) -> Unit
) : RecyclerView.Adapter<CategoriesAdapter.CategoryViewHolder>() {

    class CategoryViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        tvCategoryName = view.findViewById(R.id.tvCategoryName)
    }

    private lateinit var tvCategoryName: TextView

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): CategoryViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_category, parent, false)
        return CategoryViewHolder(view)
    }

    override fun onBindViewHolder(holder: CategoryViewHolder, position: Int) {
        val cat = categories[position]
        tvCategoryName.text = cat["name"]?.toString()

        holder.itemView.setOnClickListener {
            onCategoryClick(cat["name"]?.toString() ?: "")
        }
    }

    override fun getItemCount(): Int = categories.size

    fun setCategories(newCategories: List<Map<String, Any>>) {
        categories = newCategories
        notifyDataSetChanged()
    }
}