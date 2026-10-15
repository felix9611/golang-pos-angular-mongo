<h1>POS System</h1>
<p><strong>Golang Gin & Angular SSR | Cloud-based Web Application for POS System</strong></p>
<p>NCPL（Non-Commercial Public License）</p>

<h2>🌟 Project Overview</h2>
<p>This is a full-stack <strong>POS (Point of Sale) System</strong> built with <strong>Golang Gin Framework</strong> and <strong>Angular SSR</strong>, It is designed for businesses to manage purchase orders, transactions, inventory, and product records.</p>
<ul>
  <li>Manage purchase orders and transaction history</li>
  <li>Track inventory and product records</li>
  <li>Generate real-time reports and export data</li>
</ul>
<p>The system is <strong>cloud-based</strong>, reducing IT infrastructure costs and offering ease of use with no installation required.</p>

<p>🔗 <strong><a href="https://pos.felix9611.com/login" target="_blank">Online Demo</a></strong></p>
<p><strong>Username:</strong> <code>Demo</code><br>
<strong>Password:</strong> <code>888888</code></p>
<p><em>Note: Backend and Frontend both hosted on self-hosted Raspberry Pi5 with Cloudflare</em></p>

<!-- Add example images here if available -->
<h2>📸 Example Screenshots</h2>

<h3>Dashboard</h3>
<img src="https://github.com/felix9611/nest-pos-mongo-angular/blob/main/image/pos-1.png" >

<h3>Stock Take Form View</h3>
<img src="https://github.com/felix9611/nest-pos-mongo-angular/blob/main/image/pos-2.png" >

<h3>Stock In</h3>
<img src="https://github.com/felix9611/nest-pos-mongo-angular/blob/main/image/pos-3.png" >

<h3>Product Form</h3>
<img src="https://github.com/felix9611/nest-pos-mongo-angular/blob/main/image/pos-4.png" >

<h3>Product List</h3>
<img src="https://github.com/felix9611/nest-pos-mongo-angular/blob/main/image/pos-5.png" >

<h3>POS</h3>
<img src="https://github.com/felix9611/nest-pos-mongo-angular/blob/main/image/pos-6.png" 

<h2>🚀 Key Features</h2>
<ul>
  <li><strong>Comprehensive Purchase Order & Transaction Management</strong></li>
  <li><strong>Real-Time Inventory Tracking</strong></li>
  <li><strong>Member System Integration</strong></li>
  <li><strong>Secure Authentication & Authorization:</strong> Utilizes <strong>JWT tokens</strong> for data protection.</li>
</ul>

<h2>⚙️ Tech Stack</h2>
<h3>Frontend:</h3>
<ul>
  <li>Angular 21</li>
  <li>Node.js 24</li>
  <li>Typescript</li>
  <li>Tailwind CSS</li>
  <li>Canvas.js 3.12.5</li>
  <!--<li>xlsx, jspdf for reporting and data export</li>-->
</ul>
<h3>Backend:</h3>
<ul>
  <li>Gin v1.10.1</li>
  <li>Mongoose 9.10 for database interaction</li>
  <li>MongoDB driver v1.17.4</li>
  <li>Gin OpenAPI UI v0.24.1 for API documentation</li>
</ul>


<h2>📋 How to Run the Project</h2>

<h3>Backend</h3>
<pre><code>// Go to the backend's file directory
cd backend

// Init library
go install 

// Dev Run
go run main.go

// Build
go build

// Clean Cache
go clean -cache

// Make API doc
swag init


// API Documentation URL
http://localhost:6600/api
</code></pre>

<h3>Frontend</h3>
<pre><code>// Go to the frontend's file directory
cd frontend

// Install dependencies (npm)
npm install
// or
yarn

// Run the frontend
npm run start
// or
yarn start

// Build the frontend
npm run build:uat or build:prod
// or
yarn build:uat or build:prod

// Preview URL
http://localhost:4200
</code></pre>

<h2>🌐 Deployment</h2>
<ul>
  <li><strong>Database:</strong> <del>Mongodb Atlas</del><strong>New experiment!</strong> Self-hosted in My Raspberry Pi5 server</li>
  <li><strong>Backend:</strong> <del>AWS runing in Linux & Nginx</del> <strong>New experiment!</strong> Fullset self-hosted server power by my Raspberry Pi5 server</li>
  <li><strong>Frontend:</strong> Cloudflare Tunnel DNS + Raspberry Pi5 runing in Angualr SSR</li>
</ul>

<h2>📈 Business Impact</h2>
<ul>
  <li><strong>Improved Business Operations</strong> with real-time inventory tracking and reporting.</li>
  <li><strong>Enhanced Data Security</strong> through JWT-based user authentication.</li>
  <li><strong>Reduced Costs</strong> by eliminating the need for on-premise installations.</li>
</ul>

<h2>Version Update Note</h2>
<ul>
  <li><strong>Node updated to version 22 at 24 Sep 2026</strong></li>
  <li><strong>Please remove package lock files and folder node_modules both Backend & Frontend folder, and reinstall node_modules</li>
</ul>