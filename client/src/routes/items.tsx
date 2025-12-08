import { useMutation, useQuery } from "@connectrpc/connect-query";
import { createFileRoute } from "@tanstack/react-router";
import { Pencil, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { Button, Card, Dialog, Input, Textarea } from "@/components/ui";
import type { Item } from "@/proto-generated/item/v1/item_pb";
import {
	createItem,
	deleteItem,
	listItems,
	updateItem,
} from "@/proto-generated/item/v1/item_service-ItemService_connectquery";

export const Route = createFileRoute("/items")({
	component: Items,
});

function Items() {
	const { data, isLoading, refetch } = useQuery(listItems, {});

	const updateMutation = useMutation(updateItem);

	const createMutation = useMutation(createItem);

	const deleteMutation = useMutation(deleteItem);

	const [showDialog, setShowDialog] = useState(false);
	const [editingItem, setEditingItem] = useState<Item | null>(null);
	const [formData, setFormData] = useState({ name: "", description: "" });

	const items = data?.items || [];

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		try {
			if (editingItem) {
				await updateMutation.mutateAsync({
					id: editingItem.id,
					name: formData.name,
					description: formData.description,
				});
			} else {
				await createMutation.mutateAsync({
					name: formData.name,
					description: formData.description,
				});
			}
			setFormData({ name: "", description: "" });
			setEditingItem(null);
			setShowDialog(false);
			refetch();
		} catch (error) {
			console.error("Failed to save item:", error);
		}
	};

	const handleEdit = (item: Item) => {
		setEditingItem(item);
		setFormData({ name: item.name, description: item.description });
		setShowDialog(true);
	};

	const handleDelete = async (id: string) => {
		if (!confirm("Are you sure you want to delete this item?")) return;
		try {
			await deleteMutation.mutateAsync({ id });
			refetch();
		} catch (error) {
			console.error("Failed to delete item:", error);
		}
	};

	const handleCancel = () => {
		setFormData({ name: "", description: "" });
		setEditingItem(null);
		setShowDialog(false);
	};

	const openCreateDialog = () => {
		setEditingItem(null);
		setFormData({ name: "", description: "" });
		setShowDialog(true);
	};

	return (
		<div className="container mx-auto px-4 py-8 max-w-6xl">
			<div className="flex justify-between items-center mb-8">
				<div>
					<h1 className="text-4xl font-bold tracking-tight">Items</h1>
					<p className="text-muted-foreground mt-2">
						Manage your items with full CRUD operations
					</p>
				</div>
				<Button onClick={openCreateDialog}>
					<Plus className="h-4 w-4 mr-2" />
					Add Item
				</Button>
			</div>

			{/* Dialog for Create/Edit */}
			<Dialog open={showDialog} onOpenChange={setShowDialog}>
				<Dialog.Content>
					<Dialog.Header>
						<Dialog.Title>
							{editingItem ? "Edit Item" : "Create New Item"}
						</Dialog.Title>
						<Dialog.Description>
							{editingItem
								? "Update the item details below"
								: "Fill in the details for your new item"}
						</Dialog.Description>
					</Dialog.Header>

					<form onSubmit={handleSubmit}>
						<div className="space-y-4 py-4">
							<Input.Field
								label="Name"
								type="text"
								value={formData.name}
								onChange={(e) =>
									setFormData({ ...formData, name: e.target.value })
								}
								required
								placeholder="Enter item name"
							/>

							<Textarea.Field
								label="Description"
								value={formData.description}
								onChange={(e) =>
									setFormData({ ...formData, description: e.target.value })
								}
								required
								rows={4}
								placeholder="Enter item description"
							/>
						</div>

						<Dialog.Footer>
							<Button type="button" variant="outline" onClick={handleCancel}>
								Cancel
							</Button>
							<Button type="submit">{editingItem ? "Update" : "Create"}</Button>
						</Dialog.Footer>
					</form>
				</Dialog.Content>
			</Dialog>

			{/* Items List */}
			{isLoading ? (
				<div className="flex items-center justify-center py-12">
					<div className="text-muted-foreground">Loading...</div>
				</div>
			) : items.length === 0 ? (
				<Card>
					<Card.Content className="flex flex-col items-center justify-center py-12">
						<p className="text-muted-foreground text-center">
							No items yet. Create your first item!
						</p>
						<Button onClick={openCreateDialog} className="mt-4">
							<Plus className="h-4 w-4 mr-2" />
							Create First Item
						</Button>
					</Card.Content>
				</Card>
			) : (
				<div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
					{items.map((item) => (
						<Card key={item.id}>
							<Card.Header>
								<Card.Title>{item.name}</Card.Title>
								<Card.Description>{item.description}</Card.Description>
							</Card.Header>

							<Card.Content>
								<div className="text-xs text-muted-foreground space-y-1">
									<p>
										Created:{" "}
										{item.createdAt
											? new Date(
													Number(item.createdAt.seconds) * 1000,
												).toLocaleDateString()
											: "N/A"}
									</p>
									{item.updatedAt && (
										<p>
											Updated:{" "}
											{new Date(
												Number(item.updatedAt.seconds) * 1000,
											).toLocaleDateString()}
										</p>
									)}
								</div>
							</Card.Content>

							<Card.Footer className="gap-2">
								<Button
									variant="outline"
									size="sm"
									onClick={() => handleEdit(item)}
									className="flex-1"
								>
									<Pencil className="h-3 w-3 mr-1" />
									Edit
								</Button>
								<Button
									variant="destructive"
									size="sm"
									onClick={() => handleDelete(item.id)}
									className="flex-1"
								>
									<Trash2 className="h-3 w-3 mr-1" />
									Delete
								</Button>
							</Card.Footer>
						</Card>
					))}
				</div>
			)}
		</div>
	);
}
