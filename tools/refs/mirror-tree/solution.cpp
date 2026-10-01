#include <utility>
using namespace std;

TreeNode* mirror(TreeNode* root) {
    if (!root) return nullptr;
    swap(root->left, root->right);
    mirror(root->left);
    mirror(root->right);
    return root;
}
