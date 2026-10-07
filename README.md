
<h1>Fixed Asset Management System In Golang</h1>
<p><strong>Golang Gin & Angular SSR | Cloud-based Web Application for Asset Management</strong></p>
<p>Current version: V1.0 </p>
<p>Next version function: Unit test and Excel Data Export </p>
<p>If you find this project helpful, please click ⭐ Star! This helps more people discover it.</p>

<h2>🌟 Project Overview</h2>
<p>This is a full-stack <strong>Fixed Asset Management System</strong> built with <strong>Golang Gin Framework</strong> and <strong>Angular SSR</strong>, designed for managing fixed assets within facilities. The system includes:</p>
<ul>
  <li>Purchase and maintenance record management</li>
  <li>Asset write-off tracking</li>
  <li>Status monitoring and reporting</li>
</ul>
<p>The system is <strong>cloud-based</strong>, significantly reducing IT infrastructure costs and improving usability, with no installation required.</p>

<p>🔗 <strong><a href="https://fixedasset-gin.felix9611.com/login" target="_blank">Online Demo Front from Cloudflare DNS Tunnel and my Raspberry Pi5 server</a></strong></p>

<p><strong>Username:</strong> <code>Demo</code><br>
<strong>Password:</strong> <code>888888</code></p>
<p><em>Note: Backend and Frontend both hosted on self-hosted Raspberry Pi5 with Cloudflare</em></p>

<h2>📸 Example Screenshots</h2>

<h3>Tax Information</h3>
<img src="https://github.com/felix9611/nest-fixedasset-mongo-angular/blob/dev/image/fixedasset-1.png" alt="Dashboard Overview">

<h3>Asset Form View</h3>
<img src="https://github.com/felix9611/nest-fixedasset-mongo-angular/blob/dev/image/fixedasset-2.png" alt="Asset List View">

<h3>User Info</h3>
<img src="https://github.com/felix9611/nest-fixedasset-mongo-angular/blob/dev/image/fixedasset-3.png" alt="Maintenance Records">

<h3>Dashboard</h3>
<img src="https://github.com/felix9611/nest-fixedasset-mongo-angular/blob/dev/image/fixedasset-4.png" alt="Write-Off Management">

<h2>🚀 Key Features</h2>
<ul>
  <li><strong>Comprehensive Asset Management:</strong> Manage asset status, purchase records, maintenance schedules, and write-offs.</li>
  <li><strong>Secure Authentication & Authorization:</strong> Utilizes <strong>JWT tokens</strong> for enhanced security and data protection.</li>
  <li><strong>Cloud-Based Architecture:</strong> Easy deployment and access without installation, lowering operational costs.</li>
  <li><strong>Real-Time Data Visualization:</strong> Integrated <strong>Chart.js</strong> for graphical representation of asset data.</li>
</ul>

<h2>⚙️ Tech Stack</h2>
<h3>Frontend:</h3>
<ul>
  <li>Angular 21</li>
  <li>Node.js 22</li>
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

<h2>🛠️ Core Functionalities</h2>
<ul>
  <li><strong>Asset Registration & Tracking:</strong> Record and track fixed assets with detailed status and information.</li>
  <li><strong>Expenditure & Maintenance Management:</strong> Log asset purchase, maintenance, and write-off records for better financial control.</li>
  <li><strong>Dynamic Reporting:</strong> Showing chart graphs for easy analysis and sharing.</li>
  <li><strong>Secure Authentication & RBAC:</strong> Unit by per role with menu page permission, customization user group(s) to actionable functions</li>
  <li><strong>Excel Data Import:</strong> Batch upload data increases efficiency</li>
  <li><strong>Excel Data Export:</strong> Auto batch to report data increases efficiency</li>
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
http://localhost:6500/api
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
  <li><strong>Database:</strong><strong>New experiment!</strong> Self-hosted in My Raspberry Pi5 server</li>
  <li><strong>Backend:</strong><strong>New experiment!</strong> Fullset self-hosted server power by my Raspberry Pi5 server</li>
  <li><strong>Frontend:</strong>Cloudflare Tunnel DNS + Raspberry Pi5 runing in Angualr SSR</li>
</ul>

<h2>📈 Business Impact</h2>
<ul>
  <li><strong>Reduced Operational Costs</strong> by eliminating the need for desktop installations.</li>
  <li><strong>Improved Data Security</strong> through JWT-based user authentication.</li>
  <li><strong>Enhanced Decision-Making</strong> with real-time data visualization and comprehensive reporting.</li>
</ul>
