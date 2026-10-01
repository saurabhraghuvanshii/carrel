static bool same(TreeNode* a, TreeNode* b) {
    if (!a || !b) return a == b;
    return a->val == b->val && same(a->left, b->left) && same(a->right, b->right);
}

bool hasSubtree(TreeNode* root, TreeNode* sub) {
    if (!root) return false;
    return same(root, sub) || hasSubtree(root->left, sub) || hasSubtree(root->right, sub);
}
