import React, { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [items, setItems] = useState([]);
  const [cart, setCart] = useState([]);
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");
  const [loading, setLoading] = useState(true);
  const [orderLoading, setOrderLoading] = useState(false);

  useEffect(() => {
    // We fetch from the API Gateway's route, not the service directly
    setLoading(true);
    fetch("/api/catalog/items")
      .then((res) => res.json())
      .then((data) => {
        setItems(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching items:", err);
        setMessage("Failed to load menu items. Please refresh the page.");
        setMessageType("error");
        setLoading(false);
      });
  }, []);

  const addToCart = (item) => {
    setCart((prevCart) => [...prevCart, item]);
    setMessage(`${item.name} added to cart!`);
    setMessageType("success");
    setTimeout(() => setMessage(""), 3000);
  };

  const removeFromCart = (indexToRemove) => {
    setCart((prevCart) =>
      prevCart.filter((_, index) => index !== indexToRemove)
    );
  };

  const getTotalPrice = () => {
    return cart.reduce((total, item) => total + item.price, 0);
  };

  const placeOrder = () => {
    if (cart.length === 0) {
      setMessage("Your cart is empty!");
      setMessageType("error");
      setTimeout(() => setMessage(""), 3000);
      return;
    }

    setOrderLoading(true);
    const order = {
      item_ids: cart.map((item) => item.id),
    };

    fetch("/api/orders/orders", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(order),
    })
      .then((res) => res.json())
      .then((data) => {
        setMessage(
          `Order ${
            data.id
          } placed successfully! Total: $${getTotalPrice().toFixed(2)}`
        );
        setMessageType("success");
        setCart([]); // Clear cart
        setOrderLoading(false);
        setTimeout(() => setMessage(""), 5000);
      })
      .catch((err) => {
        setMessage("Failed to place order. Please try again.");
        setMessageType("error");
        setOrderLoading(false);
        console.error("Error placing order:", err);
        setTimeout(() => setMessage(""), 5000);
      });
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>🍕 Student Cafe ☕</h1>
      </header>
      <main className="container">
        <div className="menu">
          <h2>📋 Our Menu</h2>
          {loading ? (
            <div
              style={{ textAlign: "center", padding: "2rem", color: "#7f8c8d" }}
            >
              Loading delicious items... 🍽️
            </div>
          ) : (
            <ul>
              {items.map((item) => (
                <li key={item.id}>
                  <span>
                    {item.name} - <strong>${item.price.toFixed(2)}</strong>
                  </span>
                  <button onClick={() => addToCart(item)}>
                    Add to Cart 🛒
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
        <div className="cart">
          <h2>🛒 Your Cart ({cart.length})</h2>
          {cart.length === 0 ? (
            <div className="empty-cart">
              Your cart is empty 🥺
              <br />
              Add some delicious items!
            </div>
          ) : (
            <>
              <ul>
                {cart.map((item, index) => (
                  <li key={index}>
                    <div
                      style={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                      }}
                    >
                      <span>
                        {item.name} - ${item.price.toFixed(2)}
                      </span>
                      <button
                        onClick={() => removeFromCart(index)}
                        style={{
                          background: "#e74c3c",
                          color: "white",
                          border: "none",
                          borderRadius: "50%",
                          width: "25px",
                          height: "25px",
                          cursor: "pointer",
                          fontSize: "12px",
                        }}
                      >
                        ✕
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
              <div className="cart-summary">
                <div className="cart-total">
                  Total: ${getTotalPrice().toFixed(2)}
                </div>
              </div>
            </>
          )}
          <button
            onClick={placeOrder}
            className="order-btn"
            disabled={cart.length === 0 || orderLoading}
          >
            {orderLoading
              ? "Placing Order... ⏳"
              : `Place Order (${cart.length} items) 🚀`}
          </button>
          {message && (
            <div className={`message ${messageType}`}>
              {messageType === "success" ? "✅ " : "❌ "}
              {message}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;
