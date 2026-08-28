package com.expensetracker.family.ui.profile

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import androidx.recyclerview.widget.RecyclerView
import com.expensetracker.family.R

class ProfileTransactionAdapter(
    private var transactions: List<Map<String, Any>>
) : RecyclerView.Adapter<ProfileTransactionAdapter.TransactionViewHolder>() {

    class TransactionViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        tvType = view.findViewById(R.id.tvType)
        tvAmount = view.findViewById(R.id.tvAmount)
        tvDate = view.findViewById(R.id.tvDate)
        tvDescription = view.findViewById(R.id.tvDescription)
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): TransactionViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_profile_transaction, parent, false)
        return TransactionViewHolder(view)
    }

    override fun onBindViewHolder(holder: TransactionViewHolder, position: Int) {
        val tx = transactions[position]
        holder.tvType.text = tx["type"]?.toString()
        holder.tvAmount.text = "${tx["amount"]} ₹"
        holder.tvDate.text = tx["date"]?.toString() ?: "-"
        holder.tvDescription.text = tx["description"]?.toString() ?: "-"
    }

    override fun getItemCount(): Int = transactions.size

    fun setTransactions(newTransactions: List<Map<String, Any>>) {
        transactions = newTransactions
        notifyDataSetChanged()
    }
}