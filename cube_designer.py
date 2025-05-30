#!/usr/bin/env python3

import sys
import subprocess
import os
from PyQt5.QtWidgets import (QApplication, QMainWindow, QWidget, QVBoxLayout, 
                             QHBoxLayout, QGridLayout, QPushButton, QLabel, 
                             QTextEdit, QFrame, QLineEdit, QMessageBox, QListWidget,
                             QListWidgetItem, QProgressBar, QSplitter)
from PyQt5.QtCore import Qt, QThread, pyqtSignal
from PyQt5.QtGui import QPalette

class SolverThread(QThread):
    solution_ready = pyqtSignal(list)  # List of (step, move, cube_string) tuples
    error_occurred = pyqtSignal(str)
    
    def __init__(self, cube_string, working_dir):
        super().__init__()
        self.cube_string = cube_string
        self.working_dir = working_dir
    
    def run(self):
        try:
            result = subprocess.run(
                ['go', 'run', '.', self.cube_string],
                cwd=self.working_dir,
                capture_output=True,
                text=True,
                timeout=60
            )
            
            if result.returncode == 0:
                # Parse the solution steps
                solution_steps = []
                lines = result.stdout.strip().split('\n')
                
                for line in lines:
                    if '|' in line and line.startswith('step'):
                        parts = line.split(': ', 1)
                        if len(parts) == 2:
                            step_num = parts[0].replace('step ', '')
                            move_and_cube = parts[1].split('|')
                            if len(move_and_cube) == 2:
                                move, cube_string = move_and_cube
                                solution_steps.append((step_num, move, cube_string))
                
                self.solution_ready.emit(solution_steps)
            else:
                error_msg = result.stderr.strip() if result.stderr else "Unknown error occurred"
                self.error_occurred.emit(f"Failed to solve cube:\n{error_msg}")
                
        except subprocess.TimeoutExpired:
            self.error_occurred.emit("Solver timed out after 60 seconds. The cube might be too complex.")
        except FileNotFoundError:
            self.error_occurred.emit("Go is not installed or not in PATH. Please install Go to use the solver.")
        except Exception as e:
            self.error_occurred.emit(f"An error occurred: {str(e)}")

class CubeSquare(QPushButton):
    def __init__(self, face_idx, row, col):
        super().__init__()
        self.face_idx = face_idx
        self.row = row
        self.col = col
        self.color_value = 0  # Default to green
        self.setFixedSize(40, 40)
        self.update_color()
        
    def update_color(self):
        colors = {
            0: '#00AA00',  # green
            1: '#FFFF00',  # yellow
            2: '#FFFFFF',  # white
            3: '#FF0000',  # red
            4: '#FF8800',  # orange
            5: '#0000FF'   # blue
        }
        color = colors[self.color_value]
        self.setStyleSheet(f"""
            QPushButton {{
                background-color: {color};
                border: 2px solid #333333;
                border-radius: 3px;
            }}
            QPushButton:hover {{
                border: 2px solid #666666;
            }}
        """)

class CubeDesigner(QMainWindow):
    def __init__(self):
        super().__init__()
        self.selected_color = 0
        self.squares = {}  # Store all cube squares
        self.init_ui()
        
    def init_ui(self):
        self.setWindowTitle('Rubik\'s Cube Designer')
        self.setGeometry(100, 100, 1200, 700)
        
        central_widget = QWidget()
        self.setCentralWidget(central_widget)
        
        # Create horizontal splitter for main layout
        main_splitter = QSplitter(Qt.Horizontal)
        main_layout = QHBoxLayout(central_widget)
        main_layout.addWidget(main_splitter)
        
        # Left side - cube designer
        left_widget = QWidget()
        left_layout = QVBoxLayout(left_widget)
        
        # Color palette
        color_frame = QFrame()
        color_frame.setFrameStyle(QFrame.Box)
        color_layout = QHBoxLayout(color_frame)
        color_layout.addWidget(QLabel("Select Color:"))
        
        self.color_buttons = []
        colors = [
            ('Green', 0, '#00AA00'),
            ('Yellow', 1, '#FFFF00'),
            ('White', 2, '#FFFFFF'),
            ('Red', 3, '#FF0000'),
            ('Orange', 4, '#FF8800'),
            ('Blue', 5, '#0000FF')
        ]
        
        for name, value, color in colors:
            btn = QPushButton(name)
            btn.setFixedSize(80, 30)
            btn.setStyleSheet(f"""
                QPushButton {{
                    background-color: {color};
                    border: 2px solid #333333;
                    color: {'black' if color == '#FFFF00' or color == '#FFFFFF' else 'white'};
                    font-weight: bold;
                }}
                QPushButton:checked {{
                    border: 3px solid #000000;
                }}
            """)
            btn.setCheckable(True)
            btn.clicked.connect(lambda checked, v=value: self.select_color(v))
            self.color_buttons.append(btn)
            color_layout.addWidget(btn)
            
        self.color_buttons[0].setChecked(True)  # Green selected by default
        color_layout.addStretch()
        left_layout.addWidget(color_frame)
        
        # Cube layout (unfolded)
        cube_widget = QWidget()
        cube_layout = QGridLayout(cube_widget)
        cube_layout.setSpacing(5)
        
        # Face positions in unfolded layout:
        #     [3]      (Top)
        # [2] [0] [1] [5]  (Left, Front, Right, Back)
        #     [4]      (Bottom)
        
        face_positions = {
            3: (0, 1),  # Top
            2: (1, 0),  # Left  
            0: (1, 1),  # Front
            1: (1, 2),  # Right
            5: (1, 3),  # Back
            4: (2, 1)   # Bottom
        }
        
        face_names = {
            0: 'Front', 1: 'Right', 2: 'Left', 
            3: 'Top', 4: 'Bottom', 5: 'Back'
        }
        
        for face_idx, (grid_row, grid_col) in face_positions.items():
            face_widget = QWidget()
            face_layout = QGridLayout(face_widget)
            face_layout.setSpacing(2)
            
            # Add face label
            label = QLabel(face_names[face_idx])
            label.setAlignment(Qt.AlignCenter)
            label.setStyleSheet("font-weight: bold; margin-bottom: 5px;")
            face_layout.addWidget(label, 0, 0, 1, 3)
            
            # Add 3x3 grid of squares
            for row in range(3):
                for col in range(3):
                    square = CubeSquare(face_idx, row, col)
                    square.clicked.connect(lambda checked, s=square: self.color_square(s))
                    face_layout.addWidget(square, row + 1, col)
                    
                    # Store square for later reference
                    self.squares[(face_idx, row, col)] = square
            
            cube_layout.addWidget(face_widget, grid_row, grid_col)
        
        left_layout.addWidget(cube_widget)
        
        # Input/Output section
        io_frame = QFrame()
        io_frame.setFrameStyle(QFrame.Box)
        io_layout = QVBoxLayout(io_frame)
        
        # Input section
        input_layout = QVBoxLayout()
        input_layout.addWidget(QLabel("Enter Cube String (54 digits 0-5):"))
        
        input_row = QHBoxLayout()
        self.input_text = QLineEdit()
        self.input_text.setPlaceholderText("Enter 54-character string (e.g., 000000000111111111222222222333333333444444444555555555)")
        input_row.addWidget(self.input_text)
        
        load_btn = QPushButton("Load String")
        load_btn.clicked.connect(self.load_string)
        input_row.addWidget(load_btn)
        
        input_layout.addLayout(input_row)
        io_layout.addLayout(input_layout)
        
        # Output section
        output_layout = QVBoxLayout()
        output_layout.addWidget(QLabel("Current Cube String:"))
        self.output_text = QTextEdit()
        self.output_text.setMaximumHeight(80)
        self.output_text.setReadOnly(True)
        output_layout.addWidget(self.output_text)
        
        # Validity status
        self.validity_label = QLabel("Valid cube")
        self.validity_label.setAlignment(Qt.AlignCenter)
        self.validity_label.setStyleSheet("color: green; font-weight: bold; font-size: 14px;")
        output_layout.addWidget(self.validity_label)
        
        io_layout.addLayout(output_layout)
        
        # Buttons
        button_layout = QHBoxLayout()
        
        generate_btn = QPushButton("Generate String")
        generate_btn.clicked.connect(self.generate_string)
        button_layout.addWidget(generate_btn)
        
        reset_btn = QPushButton("Reset Cube")
        reset_btn.clicked.connect(self.reset_cube)
        button_layout.addWidget(reset_btn)
        
        # Add solve button
        self.solve_btn = QPushButton("Solve Cube")
        self.solve_btn.clicked.connect(self.solve_cube)
        button_layout.addWidget(self.solve_btn)
        
        button_layout.addStretch()
        io_layout.addLayout(button_layout)
        
        left_layout.addWidget(io_frame)
        main_splitter.addWidget(left_widget)
        
        # Right side - solution panel
        right_widget = QWidget()
        right_layout = QVBoxLayout(right_widget)
        
        # Solution header
        solution_header = QLabel("Solution Steps")
        solution_header.setAlignment(Qt.AlignCenter)
        solution_header.setStyleSheet("font-weight: bold; font-size: 16px; margin: 10px;")
        right_layout.addWidget(solution_header)
        
        # Progress bar for loading
        self.progress_bar = QProgressBar()
        self.progress_bar.setVisible(False)
        self.progress_bar.setRange(0, 0)  # Indeterminate progress
        right_layout.addWidget(self.progress_bar)
        
        # Solution list
        self.solution_list = QListWidget()
        self.solution_list.itemClicked.connect(self.on_solution_step_clicked)
        right_layout.addWidget(self.solution_list)
        
        # Clear solution button
        clear_solution_btn = QPushButton("Clear Solution")
        clear_solution_btn.clicked.connect(self.clear_solution)
        right_layout.addWidget(clear_solution_btn)
        
        main_splitter.addWidget(right_widget)
        main_splitter.setSizes([800, 400])  # Set initial sizes
        
        # Initialize with solved cube
        self.reset_cube()
        
    def select_color(self, color_value):
        self.selected_color = color_value
        for i, btn in enumerate(self.color_buttons):
            btn.setChecked(i == color_value)
    
    def color_square(self, square):
        square.color_value = self.selected_color
        square.update_color()
        self.generate_string()
        self.check_validity()
    
    def reset_cube(self):
        # Initialize solved cube: each face has its own color
        for face_idx in range(6):
            for row in range(3):
                for col in range(3):
                    square = self.squares[(face_idx, row, col)]
                    square.color_value = face_idx
                    square.update_color()
        self.generate_string()
        self.check_validity()
    
    def generate_string(self):
        # Generate string in the format: face0 + face1 + face2 + face3 + face4 + face5
        # Each face: row0 + row1 + row2 (left to right, top to bottom)
        result = ""
        
        for face_idx in range(6):
            for row in range(3):
                for col in range(3):
                    square = self.squares[(face_idx, row, col)]
                    result += str(square.color_value)
        
        self.output_text.setText(result)
    
    def load_string(self):
        cube_string = self.input_text.text().strip()
        
        # Validate input
        if len(cube_string) != 54:
            QMessageBox.warning(self, "Invalid Input", 
                              f"Cube string must be exactly 54 characters long. Got {len(cube_string)} characters.")
            return
        
        # Check if all characters are valid digits 0-5
        for char in cube_string:
            if not char.isdigit() or int(char) not in range(6):
                QMessageBox.warning(self, "Invalid Input", 
                                  f"Cube string must contain only digits 0-5. Found invalid character: '{char}'")
                return
        
        # Load the cube state
        try:
            idx = 0
            for face_idx in range(6):
                for row in range(3):
                    for col in range(3):
                        square = self.squares[(face_idx, row, col)]
                        square.color_value = int(cube_string[idx])
                        square.update_color()
                        idx += 1
            
            # Update the output text and clear input
            self.generate_string()
            self.check_validity()
            self.input_text.clear()
            
        except Exception as e:
            QMessageBox.critical(self, "Error", f"Failed to load cube string: {str(e)}")
    
    def check_validity(self):
        # Count the number of each color
        color_counts = [0] * 6
        
        for face_idx in range(6):
            for row in range(3):
                for col in range(3):
                    square = self.squares[(face_idx, row, col)]
                    color_counts[square.color_value] += 1
        
        # Check if each color appears exactly 9 times
        is_valid = all(count == 9 for count in color_counts)
        
        if is_valid:
            self.validity_label.setText("Valid cube")
            self.validity_label.setStyleSheet("color: green; font-weight: bold; font-size: 14px;")
            self.solve_btn.setEnabled(True)
        else:
            invalid_colors = []
            for i, count in enumerate(color_counts):
                if count != 9:
                    color_names = ['Green', 'Yellow', 'White', 'Red', 'Orange', 'Blue']
                    invalid_colors.append(f"{color_names[i]}: {count}")
            
            self.validity_label.setText(f"Invalid cube - {', '.join(invalid_colors)}")
            self.validity_label.setStyleSheet("color: red; font-weight: bold; font-size: 14px;")
            self.solve_btn.setEnabled(False)
    
    def solve_cube(self):
        cube_string = self.output_text.toPlainText()
        current_dir = os.path.dirname(os.path.abspath(__file__))
        
        # Show progress bar and disable solve button
        self.progress_bar.setVisible(True)
        self.solve_btn.setEnabled(False)
        self.solution_list.clear()
        
        # Start solver thread
        self.solver_thread = SolverThread(cube_string, current_dir)
        self.solver_thread.solution_ready.connect(self.on_solution_ready)
        self.solver_thread.error_occurred.connect(self.on_solver_error)
        self.solver_thread.start()
    
    def on_solution_ready(self, solution_steps):
        # Hide progress bar and re-enable solve button
        self.progress_bar.setVisible(False)
        self.solve_btn.setEnabled(True)
        
        # Populate solution list
        for step_num, move, cube_string in solution_steps:
            item_text = f"Step {step_num}: {move}"
            item = QListWidgetItem(item_text)
            item.setData(Qt.UserRole, cube_string)  # Store cube string in item data
            self.solution_list.addItem(item)
    
    def on_solver_error(self, error_message):
        # Hide progress bar and re-enable solve button
        self.progress_bar.setVisible(False)
        self.solve_btn.setEnabled(True)
        
        # Show error message
        QMessageBox.critical(self, "Solver Error", error_message)
    
    def on_solution_step_clicked(self, item):
        # Get cube string from item data
        cube_string = item.data(Qt.UserRole)
        if cube_string:
            self.load_cube_from_string(cube_string)
    
    def load_cube_from_string(self, cube_string):
        """Load cube state from string without validation (internal use)"""
        try:
            idx = 0
            for face_idx in range(6):
                for row in range(3):
                    for col in range(3):
                        square = self.squares[(face_idx, row, col)]
                        square.color_value = int(cube_string[idx])
                        square.update_color()
                        idx += 1
            
            self.generate_string()
            self.check_validity()
            
        except Exception as e:
            QMessageBox.critical(self, "Error", f"Failed to load cube state: {str(e)}")
    
    def clear_solution(self):
        self.solution_list.clear()

def main():
    app = QApplication(sys.argv)
    window = CubeDesigner()
    window.show()
    sys.exit(app.exec_())

if __name__ == '__main__':
    main()