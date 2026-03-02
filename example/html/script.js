const API_BASE = "http://localhost:8080/api/product";

const listEl = document.getElementById("list");
const errorEl = document.getElementById("error");

const modal = document.getElementById("modal");
const modalTitle = document.getElementById("modalTitle");

const openCreateBtn = document.getElementById("openCreateBtn");
const closeModalBtn = document.getElementById("closeModalBtn");
const saveBtn = document.getElementById("saveBtn");

const titleInput = document.getElementById("title");
const descriptionInput = document.getElementById("description");
const priceInput = document.getElementById("price");
const photoInput = document.getElementById("photo");

let editingId = null; // если null — создаём, иначе редактируем

// ======================
// Загрузка списка
// ======================
async function loadProducts() {
    errorEl.textContent = "";
    listEl.innerHTML = "Загрузка...";

    try {
        const res = await fetch(API_BASE + "/all");

        if (!res.ok) {
            const errText = await res.json();
            throw new Error(errText);
        }

        const products = await res.json();
        renderProducts(products);

    } catch (err) {
        listEl.innerHTML = "";
        errorEl.textContent = err.message || "Ошибка загрузки";
    }
}

// ======================
// Отрисовка
// ======================
function renderProducts(products) {
    listEl.innerHTML = "";

    if (!products || products.length === 0) {
        listEl.innerHTML = "Товаров нет";
        return;
    }

    products.forEach(p => {
        const div = document.createElement("div");
        div.className = "item";

        const photoHtml = p.photo
            ? `<img src="${p.photo}" alt="${p.title}">`
            : "";

        div.innerHTML = `
      ${photoHtml}
      <h3>${p.title}</h3>
      <div>${p.description || ""}</div>
      <div class="price">Цена: ${p.price}</div>
      <button class="edit-btn">Изменить</button>
      <button class="delete-btn">Удалить</button>
    `;

        // Удаление
        div.querySelector(".delete-btn").addEventListener("click", async () => {
            if (!confirm("Удалить товар?")) return;

            try {
                const res = await fetch(API_BASE + "?id=" + p.id, {
                    method: "DELETE"
                });

                if (!res.ok) {
                    const errText = await res.json();
                    throw new Error(errText);
                }

                loadProducts();

            } catch (err) {
                errorEl.textContent = err.message || "Ошибка удаления";
            }
        });

        // Редактирование
        div.querySelector(".edit-btn").addEventListener("click", async () => {
            try {
                const res = await fetch(API_BASE + "?id=" + p.id);

                if (!res.ok) {
                    const errText = await res.json();
                    throw new Error(errText);
                }

                const product = await res.json();

                editingId = product.id;

                modalTitle.textContent = "Изменить товар";

                titleInput.value = product.title;
                descriptionInput.value = product.description;
                priceInput.value = product.price;
                photoInput.value = product.photo;

                modal.classList.remove("hidden");

            } catch (err) {
                errorEl.textContent = err.message || "Ошибка загрузки товара";
            }
        });

        listEl.appendChild(div);
    });
}

// ======================
// Открытие для создания
// ======================
openCreateBtn.addEventListener("click", () => {
    editingId = null;
    modalTitle.textContent = "Создать товар";

    titleInput.value = "";
    descriptionInput.value = "";
    priceInput.value = "";
    photoInput.value = "";

    modal.classList.remove("hidden");
});

// ======================
// Закрытие
// ======================
closeModalBtn.addEventListener("click", () => {
    modal.classList.add("hidden");
});

// ======================
// Сохранение (Create / Update)
// ======================
saveBtn.addEventListener("click", async () => {
    try {
        let method = "POST";
        let bodyData = {
            title: titleInput.value,
            description: descriptionInput.value,
            price: Number(priceInput.value),
            photo: photoInput.value
        };

        if (editingId) {
            method = "PUT";
            bodyData.id = editingId;
        }

        const res = await fetch(API_BASE, {
            method: method,
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(bodyData)
        });

        if (!res.ok) {
            const errText = await res.json();
            throw new Error(errText);
        }

        modal.classList.add("hidden");
        loadProducts();

    } catch (err) {
        errorEl.textContent = err.message || "Ошибка сохранения";
    }
});

// ======================
loadProducts();
