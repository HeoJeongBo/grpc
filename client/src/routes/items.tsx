import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { itemClient } from "@/lib/client";
import type { Item } from "@/proto-generated/item/v1/item_pb";

export const Route = createFileRoute("/items")({
	component: Items,
});

function Items() {
	const [items, setItems] = useState<Item[]>([]);
	const [loading, setLoading] = useState(false);
	const [showForm, setShowForm] = useState(false);
	const [editingItem, setEditingItem] = useState<Item | null>(null);
	const [formData, setFormData] = useState({ name: "", description: "" });

	const loadItems = async () => {
		setLoading(true);
		try {
			const response = await itemClient.listItems({});
			setItems(response.items || []);
		} catch (error) {
			console.error("Failed to load items:", error);
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadItems()
	}, [])

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		try {
			if (editingItem) {
				await itemClient.updateItem({
					id: editingItem.id,
					name: formData.name,
					description: formData.description,
				});
			} else {
				await itemClient.createItem({
					name: formData.name,
					description: formData.description,
				});
			}
			setFormData({ name: "", description: "" });
			setEditingItem(null);
			setShowForm(false);
			loadItems();
		} catch (error) {
			console.error("Failed to save item:", error);
		}
	};

	const handleEdit = (item: Item) => {
		setEditingItem(item);
		setFormData({ name: item.name, description: item.description });
		setShowForm(true);
	};

	const handleDelete = async (id: string) => {
		if (!confirm("Are you sure you want to delete this item?")) return;
		try {
			await itemClient.deleteItem({ id });
			loadItems();
		} catch (error) {
			console.error("Failed to delete item:", error);
		}
	};

	const handleCancel = () => {
		setFormData({ name: "", description: "" });
		setEditingItem(null);
		setShowForm(false);
	};

	return (
		<div className="px-4 py-6 sm:px-0">
			<div className="flex justify-between items-center mb-6">
				<h1 className="text-3xl font-bold text-gray-900">Items</h1>
				<button
					onClick={() => setShowForm(true)}
					className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
				>
					Add Item
				</button>
			</div>

			{showForm && (
				<div className="mb-6 bg-white p-6 rounded-lg shadow">
					<h2 className="text-xl font-semibold mb-4">
						{editingItem ? "Edit Item" : "Create New Item"}
					</h2>
					<form onSubmit={handleSubmit}>
						<div className="mb-4">
							<label className="block text-sm font-medium text-gray-700 mb-2">
								Name
							</label>
							<input
								type="text"
								value={formData.name}
								onChange={(e) =>
									setFormData({ ...formData, name: e.target.value })
								}
								className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								required
							/>
						</div>
						<div className="mb-4">
							<label className="block text-sm font-medium text-gray-700 mb-2">
								Description
							</label>
							<textarea
								value={formData.description}
								onChange={(e) =>
									setFormData({ ...formData, description: e.target.value })
								}
								className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								rows={3}
								required
							/>
						</div>
						<div className="flex gap-2">
							<button
								type="submit"
								className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
							>
								{editingItem ? "Update" : "Create"}
							</button>
							<button
								type="button"
								onClick={handleCancel}
								className="bg-gray-300 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-400"
							>
								Cancel
							</button>
						</div>
					</form>
				</div>
			)}

			{loading ? (
				<div className="text-center py-8">Loading...</div>
			) : (
				<div className="bg-white shadow overflow-hidden sm:rounded-md">
					<ul className="divide-y divide-gray-200">
						{items.map((item) => (
							<li key={item.id} className="px-6 py-4 hover:bg-gray-50">
								<div className="flex items-center justify-between">
									<div className="flex-1">
										<h3 className="text-lg font-medium text-gray-900">
											{item.name}
										</h3>
										<p className="text-sm text-gray-600 mt-1">
											{item.description}
										</p>
										<p className="text-xs text-gray-400 mt-2">
											Created:{" "}
											{item.createdAt
												? new Date(
														Number(item.createdAt.seconds) * 1000,
													).toLocaleString()
												: "N/A"}
										</p>
									</div>
									<div className="flex gap-2">
										<button
											onClick={() => handleEdit(item)}
											className="text-indigo-600 hover:text-indigo-900 px-3 py-1"
										>
											Edit
										</button>
										<button
											onClick={() => handleDelete(item.id)}
											className="text-red-600 hover:text-red-900 px-3 py-1"
										>
											Delete
										</button>
									</div>
								</div>
							</li>
						))}
						{items.length === 0 && (
							<li className="px-6 py-8 text-center text-gray-500">
								No items yet. Create your first item!
							</li>
						)}
					</ul>
				</div>
			)}
		</div>
	);
}
