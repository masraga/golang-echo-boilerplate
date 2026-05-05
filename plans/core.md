# Business Plan: Custom Window Treatment & Textile Fabrication
**Business Model:** Hybrid Retailer-Manufacturer (Value-Added Textile Center)

---

## 1. Executive Summary
This business operates as a bridge between raw material distribution and finished-goods manufacturing. It serves three distinct markets: DIY retail customers, professional interior designers, and high-volume commercial contractors. The core value proposition is the ability to transform raw fabric rolls (measured in meters) into custom-engineered window solutions using a structured Bill of Materials (BOM) logic.

---

## 2. The Core Pillars

### A. Inventory Layer
* **Metered SKUs:** Fabric rolls handled in decimals (e.g., 50.75m). Requires tracking of "Dye Lots" to ensure color consistency across orders.
* **Unit SKUs:** Discrete hardware items (hooks, rings, rods, brackets) handled in whole integers.
* **Conversion Logic:** Ability to buy in bulk (Rolls/Crates) and sell in fractions (Meters/Units).

### B. Production Layer (The Workroom)
* **Bill of Materials (BOM):** Every custom product is a "recipe." 
    * *Formula:* `1 Finished Curtain = (X meters of Fabric) + (Y units of Hardware) + (Z hours of Labor)`.
* **Waste Management:** Accounting for the 5–10% "shrinkage" or scrap factor inherent in textile cutting.

### C. Sales Layer (Multi-Tier Pricing)
* **Retail:** Premium margins for walk-in customers and one-off custom jobs.
* **Wholesale:** Discounted rates for B2B partners (designers/hotels) based on volume or account status.

---

## 3. Operational Flow (Step-by-Step)

### Step 1: Procurement & Intake
* **Action:** Receive goods into the Warehouse.
* **System Impact:** Update stock levels. Assign Batch/Lot numbers to fabric rolls.
* **Goal:** Ensure the "Digital Twin" of the warehouse matches the physical shelves.

### Step 2: Sales & Quoting
* **Action:** Clerk identifies customer type (Retail vs. Wholesale) and takes measurements.
* **System Impact:** System generates a quote using the appropriate price tier. 
* **Reservation:** Raw materials are "soft-locked" to prevent them from being sold to others during the production lead time.

### Step 3: The Work Order (Trigger)
* **Action:** Sales order is converted to a "Work Order" for the warehouse/workroom.
* **Instructions:** Provide workers with a precise "Cutting List" (e.g., "Cut 4.2m from Roll #A-102").

### Step 4: Fabrication & Assembly
* **Action:** Workers cut the fabric and pull hardware from bins.
* **Inventory Impact:** **Actual Deduction.** The raw fabric is removed from the roll balance, and hardware units are subtracted from stock.
* **Assembly:** Sewing, pleating, and finishing the product.

### Step 5: Quality Control & Fulfillment
* **Action:** Final measurement check. Customer is notified for pickup or installation.
* **System Impact:** Final invoice is cleared; "Work Order" is marked as "Closed."

---

## 4. Strategic Growth Phases

### Phase 1: Operational Excellence
Focus on **Yield per Roll**. By digitalizing the cutting process, the owner reduces "mystery waste" and ensures every meter of fabric is either sold or accounted for as scrap.

### Phase 2: B2B Scaling
Leverage the warehouse and worker capacity to take on **Contract Projects** (Hotels, Hospitals, Office Buildings). These high-volume orders require strict wholesale pricing logic and batch tracking.

### Phase 3: Digital Catalog & Showroom
Transition the physical store into a "Design Showroom" while the backend system handles complex logistics, allowing for online custom orders and sample shipping.

---

## 5. Key Success Metrics (KPIs)
1.  **Inventory Accuracy:** The variance between physical stock and digital records.
2.  **Scrap Rate:** Percentage of fabric lost to cutting errors or remnants.
3.  **Labor Efficiency:** Time taken from "Work Order Created" to "Product Finished."
4.  **Tier Contribution:** Revenue split between Retail (High Margin) and Wholesale (High Volume).

---

## 6. Risk Mitigation
* **Dye-Lot Management:** Always fulfill a single room's order from the same roll to avoid 2% color variances.
* **Dead Stock:** Identify slow-moving fabric rolls early and move them via "Remnant Sales" to maintain cash flow.
* **Workroom Bottlenecks:** Use digital Work Orders to track which workers are overloaded and reassign tasks in real-time.
